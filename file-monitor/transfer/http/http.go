package http_transfer

import (
	"drpbx/file-monitor/transfer"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	endpoint_create = "create"
	endpoint_update = "update"
)

type HTTPTransfer struct {
	url    string
	client *http.Client
}

func NewHTTPTransfer(url string) *HTTPTransfer {
	return &HTTPTransfer{
		url:    url,
		client: &http.Client{},
	}
}

var _ transfer.Transfer = (*HTTPTransfer)(nil)

func (h *HTTPTransfer) send(filePath, endpoint string) error {

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file: %s", err.Error())
		return err
	}
	defer file.Close()

	// Create a pipe to stream data
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		part, err := writer.CreateFormFile("file", filePath)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to create form file: %s", err.Error()))
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to copy file data: %s", err.Error()))
			return
		}
		writer.Close()

		fmt.Println("File uploaded successfully")
	}()

	req, err := http.NewRequest("POST", h.url+endpoint, pr)
	if err != nil {
		log.Printf("failed to create request: %s", err.Error())
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("upload failed: %s", body)
		return errors.New("failed to upload file")
	}

	return nil
}

func (h *HTTPTransfer) Update(filePath string) error {
	return h.send(filePath, endpoint_update)
}

func (h *HTTPTransfer) Create(filePath string) error {
	return h.send(filePath, endpoint_create)
}

func (h *HTTPTransfer) Delete(filename string) error {
	req, err := http.NewRequest("POST", h.url+"delete", strings.NewReader(filename))
	if err != nil {
		log.Printf("failed to create request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("upload failed: %s", body)
		return fmt.Errorf("failed to delete file: %s", string(body))
	}

	fmt.Println("File deleted successfully")
	return nil
}
