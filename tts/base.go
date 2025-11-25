package tts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type CachingTTSAgent interface {
	ClearCache()
	GetCachedAudioBuffer(key string) []byte
	SaveCachedAudioBuffer(key string, audioBuffer []byte) error
}

type hashCachingTTSAgent struct {
	mu       sync.Mutex
	cacheDir string
}

func NewHashCachingTTSAgent(cacheDir string) CachingTTSAgent {
	// make sure cacheDir exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}
	fmt.Println("Cache directory:", cacheDir)
	return &hashCachingTTSAgent{mu: sync.Mutex{}, cacheDir: cacheDir}
}
func (a *hashCachingTTSAgent) ClearCache() {
	a.mu.Lock()
	defer a.mu.Unlock()
}
func (a *hashCachingTTSAgent) GetCachedAudioBuffer(key string) []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	cacheKey := a.getCacheKey(key)
	cacheFile := filepath.Join(a.cacheDir, cacheKey)
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}
	return data
}

func (a *hashCachingTTSAgent) getCacheKey(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
func (a *hashCachingTTSAgent) SaveCachedAudioBuffer(key string, audioBuffer []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	cacheKey := a.getCacheKey(key)
	cacheFile := filepath.Join(a.cacheDir, cacheKey)
	return os.WriteFile(cacheFile, audioBuffer, 0644)
}
