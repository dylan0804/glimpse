package screenshots

import (
	"io/fs"
	"os"
	"path/filepath"
)

type DirProvider interface {
	GetHomeDir() (string, error)
	ReadDir(path string) ([]fs.DirEntry, error)
}

type Dir struct {}

func NewDirProvider() DirProvider {
	return &Dir{}
}

func (Dir) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (d *Dir) ReadDir(homeDir string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(filepath.Join(homeDir, "Desktop"))
	if err != nil {
		return nil, err
	}
	
	return entries, nil
}
