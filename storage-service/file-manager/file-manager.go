package file_manager

import (
	"io"
)

type FileManager interface {
	Create(file io.ReadCloser, filename string) error
	Update(file io.ReadCloser, filename string) error
	Delete(file string) error
}
