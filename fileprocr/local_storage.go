package fileprocr

import (
	"io"
	"os"
)

type LocalFileStorage struct{}

func NewLocalFileStorage() *LocalFileStorage {
	return &LocalFileStorage{}
}

func (s *LocalFileStorage) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}
