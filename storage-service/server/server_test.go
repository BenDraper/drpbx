package server

import (
	"bytes"
	fm_mocks "drpbx/storage-service/file-manager/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_write(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileManager := fm_mocks.NewMockFileManager(mockCtrl)

	tests := map[string]struct {
		mockFileManagerOutcomes func(fileManagerMocks *fm_mocks.MockFileManager)
		port                    string
		writeToFileFunc         func(io.ReadCloser, string) error
		wantStatus              int
		formField               string
		filename                string
		content                 string
	}{
		"success": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {
				fileManagerMocks.EXPECT().Create(gomock.Any(), "testName.txt").Return(nil)
			},
			port:            "8080",
			writeToFileFunc: mockFileManager.Create,
			wantStatus:      http.StatusOK,
			formField:       "file",
			filename:        "testName.txt",
			content:         "test content",
		},
		"missing file name": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {},
			port:                    "8080",
			writeToFileFunc:         mockFileManager.Create,
			wantStatus:              http.StatusBadRequest,
			formField:               "file",
			filename:                "",
			content:                 "test content",
		},
		"missing form field": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {},
			port:                    "8080",
			writeToFileFunc:         mockFileManager.Create,
			wantStatus:              http.StatusBadRequest,
			formField:               "",
			filename:                "testName.txt",
			content:                 "test content",
		},
		"write error": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {
				fileManagerMocks.EXPECT().Create(gomock.Any(), "testName.txt").Return(fmt.Errorf("test error"))
			},
			port:            "8080",
			writeToFileFunc: mockFileManager.Create,
			wantStatus:      http.StatusInternalServerError,
			formField:       "file",
			filename:        "testName.txt",
			content:         "test content",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mockFileManagerOutcomes(mockFileManager)

			s := &Server{
				fileManager: mockFileManager,
				port:        tt.port,
			}

			recorder := httptest.NewRecorder()

			req, err := createMultipartRequest(tt.formField, tt.filename, tt.content)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			s.write(recorder, req, tt.writeToFileFunc)

			resp := recorder.Result()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func createMultipartRequest(fieldName, filename, content string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if fieldName != "" {
		part, err := writer.CreateFormFile(fieldName, filename)
		if err != nil {
			return nil, err
		}
		part.Write([]byte(content))
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestServer_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileManager := fm_mocks.NewMockFileManager(mockCtrl)

	tests := map[string]struct {
		mockFileManagerOutcomes func(fileManagerMocks *fm_mocks.MockFileManager)
		port                    string
		r                       *http.Request
		wantStatus              int
	}{
		"success": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {
				fileManagerMocks.EXPECT().Delete("testName.txt").Return(nil)
			},
			port:       "8080",
			r:          createDeleteRequest("testName.txt"),
			wantStatus: http.StatusOK,
		},
		"empty body": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {

			},
			port:       "8080",
			r:          createDeleteRequest(""),
			wantStatus: http.StatusBadRequest,
		},
		"delete error": {
			mockFileManagerOutcomes: func(fileManagerMocks *fm_mocks.MockFileManager) {
				fileManagerMocks.EXPECT().Delete("testName.txt").Return(fmt.Errorf("test error"))
			},
			port:       "8080",
			r:          createDeleteRequest("testName.txt"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mockFileManagerOutcomes(mockFileManager)

			s := &Server{
				fileManager: mockFileManager,
				port:        tt.port,
			}

			recorder := httptest.NewRecorder()

			s.Delete(recorder, tt.r)

			resp := recorder.Result()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func createDeleteRequest(filename string) *http.Request {
	return httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(filename))
}
