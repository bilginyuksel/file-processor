package fileprocr_test

import (
	"os"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalFileStorageCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping local file storage integration test")
	}

	name := "local_file_storage_test.txt"

	lf := fileprocr.NewLocalFileStorage(".files")

	require.NoError(t, lf.Configure())

	wc, err := lf.Create(name)

	assert.NoError(t, err)
	assert.NotNil(t, wc)

	assert.NoError(t, os.Remove(".files/"+name))
	assert.NoError(t, os.Remove(".files/"))
}

func TestLocalFileStorageOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping local file storage integration test")
	}

	name := "local_file_storage_test.txt"

	lf := fileprocr.NewLocalFileStorage(".files")

	require.NoError(t, lf.Configure())

	wc, err := lf.Create(name)
	assert.NoError(t, err)
	assert.NotNil(t, wc)
	require.NoError(t, wc.Close())

	r, err := lf.Open(name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	require.NoError(t, r.Close())

	assert.NoError(t, os.Remove(".files/"+name))
	assert.NoError(t, os.Remove(".files/"))
}
