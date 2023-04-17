package fileprocr

import (
	"io"
	"os"
)

type LocalFileStorage struct {
	dir string
}

func NewLocalFileStorage() *LocalFileStorage {
	return &LocalFileStorage{dir: ".files"}
}

// Check if the directory is already there
// If the directory does not created yet create the directory
func (s *LocalFileStorage) Configure() error {
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		return os.Mkdir(s.dir, 0777)
	}

	return nil
}

func (s *LocalFileStorage) Create(name string) (io.WriteCloser, error) {
	return os.Create(s.dir + "/" + name)
}
