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

func (s *LocalFileStorage) Create(name string) (io.WriteCloser, error) {
	return os.Create(s.dir + "/" + name)
}
