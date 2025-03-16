package http_transfer

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPTransfer_Delete(t *testing.T) {
	tests := map[string]struct {
		client     *http.Client
		filename   string
		mockStatus int
		mockBody   string
		wantErr    bool
	}{
		"success": {
			client:     &http.Client{},
			filename:   "testfilename.txt",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		"server error": {
			client:     &http.Client{},
			filename:   "testfilename.txt",
			mockStatus: http.StatusInternalServerError,
			mockBody:   "internal server error",
			wantErr:    true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				body, _ := io.ReadAll(r.Body)
				assert.Equal(t, tt.filename, string(body))

				w.WriteHeader(tt.mockStatus)
			}))

			defer mockServer.Close()

			h := &HTTPTransfer{
				url:    mockServer.URL + "/",
				client: tt.client,
			}
			err := h.Delete(tt.filename)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestHTTPTransfer_send(t *testing.T) {
	tests := map[string]struct {
		client     *http.Client
		filename   string
		content    string
		mockStatus int
		wantErr    bool
	}{
		"success": {
			client:     &http.Client{},
			filename:   "testfilename.txt",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		"server error": {
			client:     &http.Client{},
			filename:   "testfilename.txt",
			mockStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.WriteHeader(tt.mockStatus)
			}))

			defer mockServer.Close()

			dir := t.TempDir()

			path := filepath.Join(dir, tt.filename)

			f, err := os.Create(filepath.Join(dir, tt.filename))
			assert.NoError(t, err)

			_, err = f.WriteString(tt.content)
			assert.NoError(t, err)

			h := &HTTPTransfer{
				url:    mockServer.URL + "/",
				client: tt.client,
			}

			err = h.send(path, "create")

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
