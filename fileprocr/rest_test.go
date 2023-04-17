package fileprocr_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileprocrRestHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping rest functional test")
	}

	testServer := startTestServer(t)
	defer testServer.Close()

	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test.txt")
	require.NoError(t, err)

	_, err = w.Write([]byte("Hello world!"))
	require.NoError(t, err)

	require.NoError(t, mw.Close())

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/files", buf)
	require.NoError(t, err)
	req.Header.Add("Content-Type", mw.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	var m map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&m))

	t.Log(m)

	filename := m["filename"]

	assert.NoError(t, err)
	assert.NoError(t, os.Remove(filename))
}

func startTestServer(t *testing.T) *httptest.Server {
	e := echo.New()

	lfs := fileprocr.NewLocalFileStorage("")
	require.NoError(t, lfs.Configure())
	svc := fileprocr.NewProcr(1024, lfs)
	handler := fileprocr.NewRestHandler(svc)

	handler.RegisterRoutes(e)

	return httptest.NewServer(e.Server.Handler)
}
