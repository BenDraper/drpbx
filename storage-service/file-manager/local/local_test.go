package local

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocal_Delete(t *testing.T) {
	tests := map[string]struct {
		filename  string
		writeFile bool
		wantErr   bool
	}{
		"success": {
			filename:  "test.txt",
			writeFile: true,
			wantErr:   false,
		},
		"file does not exist": {
			filename:  "test.txt",
			writeFile: false,
			wantErr:   true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()

			l := &Local{
				directory: dir,
			}

			filePath := filepath.Join(dir, tt.filename)

			if tt.writeFile {
				err := os.WriteFile(filePath, []byte("test data"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Fatalf("File does not exist before deletion")
				}
			}

			err := l.Delete(tt.filename)

			assert.Equal(t, tt.wantErr, err != nil)

			_, err = os.Stat(filePath)
			assert.True(t, os.IsNotExist(err))

		})
	}
}

func TestLocal_Create(t *testing.T) {

	tests := map[string]struct {
		content  string
		filename string
		wantErr  bool
	}{
		"success": {
			content:  "test data",
			filename: "test.txt",
			wantErr:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()

			l := &Local{
				directory: dir,
			}

			filePath := filepath.Join(dir, tt.filename)
			reader := strings.NewReader(tt.content)
			file := io.NopCloser(reader)

			err := l.Create(file, tt.filename)

			assert.Equal(t, tt.wantErr, err != nil)

			_, err = os.Stat(filePath)
			assert.Falsef(t, os.IsNotExist(err), "File was not created")

			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read created file: %v", err)
			}

			assert.Equal(t, tt.content, string(data))
		})
	}
}
