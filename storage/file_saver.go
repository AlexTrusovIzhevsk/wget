package storage

import (
	"os"
	"path/filepath"
)

type FileSaver interface {
	Save(path string, data []byte) error
}

type OsFileSaver struct {
	OutputDir string
}

func (s *OsFileSaver) Save(path string, data []byte) error {
	fullPath := filepath.Join(s.OutputDir, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}
