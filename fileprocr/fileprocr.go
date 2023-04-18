package fileprocr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

//go:generate mockgen -source=fileprocr.go -destination=mock/fileprocr.go -package=mock
type storage interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
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
func (p *Procr) Store(r io.Reader) (string, error) {
	filename := uuid.NewString()
	w, err := p.storage.Create(filename)
	if err != nil {
		zap.L().Error("failed to create a new file", zap.Error(err))
		return filename, err
	}
	defer w.Close()

	buf := make([]byte, p.chunkSize)
	for {
		bytesRead, err := r.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			zap.L().Error("failed to read from reader", zap.Error(err))
			return filename, err
		}

		if _, err := w.Write(buf[:bytesRead]); err != nil {
			zap.L().Error("failed to write the buffer", zap.Error(err))
			return filename, err
		}
	}

	err = w.Close()

	zap.L().Info("file stored to the store, checking if it is proccessable")

	return filename, err
}

// processFile opens the given file using the storage
// Unmarshals the file content to json if the content is json
// Removes the keys startswith vowel and increases the integer keys by 1000
// NOTE: There are no optimizations currently but according
// to the file contents we could've write our own marshaller
// validate the json and at the same time write the json to a new file
// Current decoder supports `More` and `Token` functionality but
// it is not the best implementation in our case. Because if you call the `Token`
// method it will pop the token and if it is not an arraj but an object
// decoding will fail. If you don't check the token and just decode it will decode
// the complete object so it is not the behavior we want again..
func (p *Procr) processFile(name string) error {
	r, err := p.storage.Open(name)
	if err != nil {
		zap.L().Error("failed to open file", zap.Error(err))
		return err
	}
	defer r.Close()

	m := make(map[string]any)
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		zap.L().Error("file could not be decoded to json", zap.Error(err))
		return err
	}

	zap.L().Info("file is decoded to json, starting modifications...")

	removeKeysStartswithVowel(m)
	increaseIntegerKeys(m)

	zap.L().Info("removed keys startswith vowels and increased int keys")

	wc, err := p.storage.Create(name + ".json")
	if err != nil {
		zap.L().Error("failed to create the json ext file", zap.Error(err))
		return err
	}
	defer wc.Close()

	return json.NewEncoder(wc).Encode(m)
}

var vowelRegex = regexp.MustCompile(`^[aeiouAEIOU]`)

func removeKeysStartswithVowel(m map[string]any) {
	for k, v := range m {
		if vowelRegex.MatchString(k) {
			delete(m, k)
			continue
		}

		traverseMapsAndArrays(v, removeKeysStartswithVowel)
	}
}

func increaseIntegerKeys(m map[string]any) {
	for k, v := range m {
		if i, err := strconv.ParseInt(k, 10, 64); err == nil {
			delete(m, k)
			m[fmt.Sprintf("%d", i+1000)] = v
		}

		traverseMapsAndArrays(v, increaseIntegerKeys)
	}
}

func traverseMapsAndArrays(val any, modifyMap func(m map[string]any)) {
	if arr, ok := val.([]any); ok {
		for _, item := range arr {
			traverseMapsAndArrays(item, modifyMap)
		}
	} else if m, ok := val.(map[string]any); ok {
		modifyMap(m)
	}
}
