package data

import (
	"os"
	"path/filepath"
)

// GetDataDir returns the directory where application data is stored
func GetDataDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	dir := filepath.Join(cacheDir, "foreign7")
	_ = os.MkdirAll(dir, 0755)
	return dir
}
