package fileprocr_test

import (
	"bytes"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/bilginyuksel/file-processor/fileprocr/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	content := "hello world!"
	buf := bytes.NewBufferString(content)

	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	proc := fileprocr.NewProcr(3, mockStorage)

	mw := &mockWriteCloser{}
	mockStorage.EXPECT().Create(gomock.Any()).Return(mw, nil)

	_, err := proc.Store(buf)
	assert.NoError(t, err)
	assert.Equal(t, content, string(mw.content))
	assert.Equal(t, 4, mw.writtenTimes)
}

func TestStore_FailToCreateFile_ReturnErr(t *testing.T) {
	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	proc := fileprocr.NewProcr(10, mockStorage)

	mockStorage.EXPECT().Create(gomock.Any()).Return(nil, assert.AnError)

	_, err := proc.Store(nil)
	assert.Error(t, err)
}

func TestStore_FailToWriteToWriter_ReturnErr(t *testing.T) {
	content := "hello!"
	buf := bytes.NewBufferString(content)

	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	proc := fileprocr.NewProcr(3, mockStorage)

	mw := &mockWriteCloser{failOnWrite: true}
	mockStorage.EXPECT().Create(gomock.Any()).Return(mw, nil)

	_, err := proc.Store(buf)
	assert.Error(t, err)
}

type mockReadCloser struct {
	*bytes.Buffer
}

func (mr *mockReadCloser) Close() error {
	return nil
}

type mockWriteCloser struct {
	writtenTimes int
	content      []byte
	failOnWrite  bool
}

func (mw *mockWriteCloser) Write(b []byte) (int, error) {
	if mw.failOnWrite {
		return 0, assert.AnError
	}

	mw.writtenTimes++
	mw.content = append(mw.content, b...)
	return 0, nil
}

func (mw *mockWriteCloser) Close() error {
	return nil
}
