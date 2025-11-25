package services

import (
	"archive/zip"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/go-chi/chi/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ModelService struct {
	cfg *config.Config
	*chi.Mux
	logger *slog.Logger
}

func NewModelService(cfg *config.Config, logger *slog.Logger) *ModelService {
	router := chi.NewRouter()
	s := &ModelService{
		cfg:    cfg,
		Mux:    router,
		logger: logger,
	}
	s.setupRoutes()
	return s
}

func (s *ModelService) setupRoutes() {
	fs := http.FileServer(http.Dir(s.cfg.ModelsConfig.ModelDir))
	s.Handle("/*", http.StripPrefix("/", fs))
	s.Post("/upload-model", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(64 << 20); err != nil { // 64MB max memory buffer
			http.Error(w, err.Error(), 400)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()
		s.logger.Info("Content-Type", "ct", r.Header.Get("Content-Type"))
		s.logger.Info("Content-Length", "cl", r.Header.Get("Content-Length"))

		//  make tmp file
		tmpFile, err := os.CreateTemp(os.TempDir(), "model-*.zip")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer tmpFile.Close()
		written, err := io.Copy(tmpFile, file)
		s.logger.Info("Written", "written", written)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := s.HandleUploadZipModelFromFile(tmpFile.Name()); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("OK"))
	})
}

func (h *ModelService) HandleUploadZipModelFromFile(filePath string) error {
	log := h.logger
	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		log.Error("cannot open zip file", "error", err, "filePath", filePath)
		return err
	}
	defer zipFile.Close()
	return h.handleUploadZipModel(&zipFile.Reader)
}

func (h *ModelService) handleUploadZipModel(zipFile *zip.Reader) error {
	log := h.logger
	var targetFolder string
	targetRoot := h.cfg.ModelsConfig.ModelDir
	for _, f := range zipFile.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".model3.json") {
			// Lấy thư mục cha
			targetFolder = filepath.Dir(f.Name)
			log.Info("Found folder:", "folder", targetFolder)
			break
		}
	}
	if targetFolder == "" {
		return errors.New("not a valid model zip file")
	}
	folderName := filepath.Base(targetFolder)
	destPath := filepath.Join(targetRoot, folderName)
	for _, f := range zipFile.File {
		if !strings.HasPrefix(f.Name, targetFolder) {
			continue
		}

		relPath := strings.TrimPrefix(f.Name, targetFolder)
		relPath = strings.TrimPrefix(relPath, "/")
		outPath := filepath.Join(destPath, relPath)

		// Thư mục
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, os.ModePerm)
			continue
		}

		// Tạo thư mục chứa file
		os.MkdirAll(filepath.Dir(outPath), os.ModePerm)

		// Mở file zip
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Tạo file dest
		dst, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(dst, rc)

		dst.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	log.Info("Extracted to:", "path", destPath)
	return nil
}

func (s *ModelService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *ModelService) ChooseModel() error {
	path, err := application.OpenFileDialog().
		SetTitle("Choose model zip file").
		AddFilter("zip filter", "*.zip:").
		PromptForSingleSelection()
	if err != nil {
		return err
	}
	return s.HandleUploadZipModelFromFile(path)

}

func (s *ModelService) DeleteModel(modelName string) error {
	modelPath := filepath.Join(s.cfg.ModelsConfig.ModelDir, modelName)
	return os.RemoveAll(modelPath)
}

type Model struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	IsActive bool   `json:"is_active"`
}

// UploadModel uploads a model from a file path (for Wails binding)
func (h *ModelService) UploadModel(filePath string) error {
	log := h.logger
	log.Info("Uploading model from file", "filePath", filePath)

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		return errors.New("only .zip files are supported")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// Use existing handler
	return h.HandleUploadZipModelFromFile(filePath)
}

func (h *ModelService) GetModelList() ([]Model, error) {
	log := h.logger
	modelDir := h.cfg.ModelsConfig.ModelDir

	entries, err := os.ReadDir(modelDir)
	if err != nil {
		return nil, err
	}

	var results []Model

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		folderPath := filepath.Join(modelDir, entry.Name())
		files, err := os.ReadDir(folderPath)

		if err != nil {
			log.Warn("cannot read folder", "folder", folderPath, "err", err)
			continue
		}

		for _, f := range files {
			if strings.HasSuffix(strings.ToLower(f.Name()), ".model3.json") {
				size, _ := getFolderSize(folderPath)
				results = append(results, Model{
					ID:       entry.Name(),
					Name:     entry.Name(),
					Path:     h.createUrlPath(filepath.Join(folderPath, f.Name())),
					Size:     size,
					IsActive: strings.EqualFold(h.cfg.CharacterConfig.Live2DModelName, entry.Name()),
				})
				break
			}
		}
	}

	return results, nil
}

func (h *ModelService) createUrlPath(fullPath string) string {
	base := filepath.Clean(h.cfg.ModelsConfig.ModelDir)
	target := filepath.Clean(fullPath)

	// Tạo đường dẫn tương đối
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return ""
	}

	// Đổi sang URL path (luôn dùng /)
	urlPath := "/" + filepath.ToSlash(rel)

	return urlPath
}

func getFolderSize(folderPath string) (int64, error) {
	var size int64

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		path := filepath.Join(folderPath, file.Name())

		if file.IsDir() {
			subSize, err := getFolderSize(path)
			if err != nil {
				return 0, err
			}
			size += subSize
			continue
		}

		info, err := file.Info()
		if err != nil {
			return 0, err
		}

		size += info.Size()
	}

	return size, nil
}
