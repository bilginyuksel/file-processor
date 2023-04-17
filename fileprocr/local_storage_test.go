package fileprocr_test

import (
	"os"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalFileStorage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping local file storage integration test")
	}

	name := "local_file_storage_test.txt"

	lf := fileprocr.NewLocalFileStorage()

	require.NoError(t, lf.Configure())

	wc, err := lf.Create(name)

	assert.NoError(t, err)
	assert.NotNil(t, wc)

	assert.NoError(t, os.Remove(".files/"+name))
	assert.NoError(t, os.Remove(".files/"))
}
