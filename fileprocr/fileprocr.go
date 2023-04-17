package fileprocr

import (
	"errors"
	"io"
)

//go:generate mockgen -source=fileprocr.go -destination=mock/fileprocr.go -package=mock
type storage interface {
	Create(name string) (io.WriteCloser, error)
}

type Procr struct {
	chunkSize int

	storage storage
}

func NewProcr(chunkSize int, storage storage) *Procr {
	return &Procr{
		chunkSize: chunkSize,
		storage:   storage,
	}
}

// Store stores the given content
// Content could be very big so instead of loading it into
// memory in one go it writes to the disk in chunks
func (p *Procr) Store(r io.Reader) error {
	w, err := p.storage.Create("test.txt")
	if err != nil {
		return err
	}
	defer w.Close()

	buf := make([]byte, p.chunkSize)
	for {
		bytesRead, err := r.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if _, err := w.Write(buf[:bytesRead]); err != nil {
			return err
		}
	}

	return nil
}
