package fileprocr_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/bilginyuksel/file-processor/fileprocr/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	content := `{
		"b": {"s": 10, "k": 20, "e": 9},
		"12": [{"a": 5, "183": "hi"}, {"5": 10}],
		"a": {"12": {}}
	}`
	buf := bytes.NewBufferString(content)

	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	mockIDGenerator := mock.NewMockidgenerator(gomock.NewController(t))
	proc := fileprocr.NewProcr(1024, mockStorage, mockIDGenerator)

	mwRaw := &mockWriteCloser{}
	mwJSON := &mockWriteCloser{}
	mr := &mockReadCloser{bytes.NewBufferString(content)}

	mockFilename := "autogen-filename"
	mockIDGenerator.EXPECT().Generate().Return(mockFilename)
	mockStorage.EXPECT().Create(mockFilename).Return(mwRaw, nil)
	mockStorage.EXPECT().Create(mockFilename+".json").Return(mwJSON, nil)
	mockStorage.EXPECT().Open(mockFilename).Return(mr, nil)

	_, err := proc.Store(buf)
	assert.NoError(t, err)
	assert.Equal(t, content, string(mwRaw.content))

	assert.NoError(t, <-proc.ProcrResultQueue)

	var actualContent map[string]any
	require.NoError(t, json.Unmarshal(mwJSON.content, &actualContent))

	expectedContent := map[string]any{
		"1012": []any{
			map[string]any{"1183": "hi"},
			map[string]any{"1005": float64(10)},
		},
		"b": map[string]any{"k": float64(20), "s": float64(10)},
	}
	assert.Equal(t, expectedContent, actualContent)
}

func TestStore_InvalidJSON(t *testing.T) {
	content := "hello world!"
	buf := bytes.NewBufferString(content)

	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	mockIDGenerator := mock.NewMockidgenerator(gomock.NewController(t))
	proc := fileprocr.NewProcr(3, mockStorage, mockIDGenerator)

	mw := &mockWriteCloser{}
	mr := &mockReadCloser{bytes.NewBufferString(content)}

	mockFilename := "autogen-filename"
	mockIDGenerator.EXPECT().Generate().Return(mockFilename)
	mockStorage.EXPECT().Create(mockFilename).Return(mw, nil)
	mockStorage.EXPECT().Open(mockFilename).Return(mr, nil)

	_, err := proc.Store(buf)
	assert.NoError(t, err)
	assert.Equal(t, content, string(mw.content))
	assert.Equal(t, 4, mw.writtenTimes)

	assert.Error(t, <-proc.ProcrResultQueue)
}

func TestStore_FailToCreateFile_ReturnErr(t *testing.T) {
	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	mockIDGenerator := mock.NewMockidgenerator(gomock.NewController(t))
	proc := fileprocr.NewProcr(10, mockStorage, mockIDGenerator)

	mockFilename := "autogen-filename"
	mockIDGenerator.EXPECT().Generate().Return(mockFilename)
	mockStorage.EXPECT().Create(gomock.Any()).Return(nil, assert.AnError)

	_, err := proc.Store(nil)
	assert.Error(t, err)
}

func TestStore_FailToWriteToWriter_ReturnErr(t *testing.T) {
	content := "hello!"
	buf := bytes.NewBufferString(content)

	mockStorage := mock.NewMockstorage(gomock.NewController(t))
	mockIDGenerator := mock.NewMockidgenerator(gomock.NewController(t))
	proc := fileprocr.NewProcr(3, mockStorage, mockIDGenerator)

	mw := &mockWriteCloser{failOnWrite: true}
	mockFilename := "autogen-filename"
	mockStorage.EXPECT().Create(gomock.Any()).Return(mw, nil)
	mockIDGenerator.EXPECT().Generate().Return(mockFilename)

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
