package file_manager

import (
	"mime/multipart"
)

type FileManager interface {
	Create(file multipart.File, filename string) error
	Update(file multipart.File, filename string) error
	Delete(file string) error
}
