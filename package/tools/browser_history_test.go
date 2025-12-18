package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetRecent(t *testing.T) {
	tmpDir := "C:/Users/BaoBao/AppData/Local/BraveSoftware/Brave-Browser/User Data/Default"
	tool, err := NewBrowserHistoryTool(tmpDir)
	if err != nil {
		t.Fatalf("init tool: %v", err)
	}

	list, err := tool.GetRecent(context.Background(), GetRecentInput{Limit: 20})
	if err != nil {
		t.Fatalf("GetRecent error: %v", err)
	}

	if len(list) != 20 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
	t.Log(list)
}

func TestGetByDomain(t *testing.T) {
	tmpDir := "C:/Users/BaoBao/AppData/Local/BraveSoftware/Brave-Browser/User Data/Default"
	tool, err := NewBrowserHistoryTool(tmpDir)
	if err != nil {
		t.Fatalf("init tool: %v", err)
	}

	list, err := tool.GetByDomain(context.Background(), GetByDomainInput{Domain: "github", Limit: 5})
	if err != nil {
		t.Fatalf("GetByDomain error: %v", err)
	}

	t.Log(list)
}

func TestTemporaryCopyDeleted(t *testing.T) {
	tmpDir := "C:/Users/BaoBao/AppData/Local/BraveSoftware/Brave-Browser/User Data/Default"
	tool, err := NewBrowserHistoryTool(tmpDir)
	if err != nil {
		t.Fatalf("init tool: %v", err)
	}

	_, err = tool.GetRecent(context.Background(), GetRecentInput{Limit: 1})
	if err != nil {
		t.Fatalf("GetRecent error: %v", err)
	}

	// Kiểm tra file copy đã bị xóa
	files, _ := os.ReadDir(tmpDir)
	for _, f := range files {
		if filepath.Ext(f.Name()) == "_tmp_copy" {
			t.Fatalf("temporary copy not removed: %s", f.Name())
		}
	}
}
