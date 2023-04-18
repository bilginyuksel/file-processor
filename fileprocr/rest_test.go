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

func TestFileProcrRestHandlerUploadFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rest handler tests")
	}

	testCases := []struct {
		desc                        string
		content                     []byte
		extraAssertionsForAsyncTask func(string, *fileprocr.Procr)
	}{
		{
			desc:    "Given file is not a JSON file, should only save raw file",
			content: []byte("Hello world!"),
			extraAssertionsForAsyncTask: func(filename string, p *fileprocr.Procr) {
				assert.Error(t, <-p.ProcrResultQueue)
				assert.Error(t, os.Remove(filename+".json"))
			},
		},
		{
			desc:    "Given file is a JSON file, should save both raw and json files",
			content: []byte(`{"1": 5, "s": 10}`),
			extraAssertionsForAsyncTask: func(filename string, p *fileprocr.Procr) {
				assert.NoError(t, <-p.ProcrResultQueue)
				assert.NoError(t, os.Remove(filename+".json"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testServer, procr := startTestServer(t)
			defer testServer.Close()

			buf := new(bytes.Buffer)
			mw := multipart.NewWriter(buf)
			w, err := mw.CreateFormFile("file", "test.txt")
			require.NoError(t, err)

			_, err = w.Write(tC.content)
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
			tC.extraAssertionsForAsyncTask(filename, procr)
		})
	}
}

func startTestServer(t *testing.T) (*httptest.Server, *fileprocr.Procr) {
	e := echo.New()

	lfs := fileprocr.NewLocalFileStorage("")
	idgen := fileprocr.NewIDGenerator()
	require.NoError(t, lfs.Configure())
	svc := fileprocr.NewProcr(1024, lfs, idgen)
	handler := fileprocr.NewRestHandler(svc)

	handler.RegisterRoutes(e)

	return httptest.NewServer(e.Server.Handler), svc
}
