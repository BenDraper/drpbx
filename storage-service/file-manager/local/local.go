package local

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	directory string
}

func NewLocal(directory string) *Local {
	return &Local{
		directory: directory,
	}
}

func (l *Local) Create(file io.ReadCloser, filename string) error {
	defer file.Close()

	dstPath := filepath.Join(l.directory, filename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, file); err != nil {
		return err
	}

	log.Printf("Created file: %s", dstPath)
	return nil
}

func (l *Local) Update(file io.ReadCloser, filename string) error {
	if err := l.Delete(filename); err != nil {
		return err
	}

	if err := l.Create(file, filename); err != nil {
		return err
	}

	log.Printf("Updated file: %s", filename)
	return nil
}

func (l *Local) Delete(filename string) error {

	if err := os.Remove(filepath.Join(l.directory, filename)); err != nil {
		return err
	}

	log.Printf("Deleted file: %s", filename)
	return nil
}
