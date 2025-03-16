package server

import (
	file_manager "drpbx/storage-service/file-manager"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type Server struct {
	fileManager file_manager.FileManager
	port        string
}

func NewServer(fileManager file_manager.FileManager, port string) *Server {
	return &Server{
		fileManager: fileManager,
		port:        port,
	}
}

func (s *Server) Serve() {
	http.HandleFunc("/create", s.Create)
	http.HandleFunc("/update", s.Update)
	http.HandleFunc("/delete", s.Delete)

	log.Printf("Listening on port %s", s.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}

}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	log.Printf("got update request")
	s.write(w, r, s.fileManager.Update)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("got create request")
	s.write(w, r, s.fileManager.Create)
}

func (s *Server) write(w http.ResponseWriter, r *http.Request, writeToFileFunc func(multipart.File, string) error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB memory buffer
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := writeToFileFunc(file, fileHeader.Filename); err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", fileHeader.Filename)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("got delete request")

	file, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.fileManager.Delete(string(file)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
