package main

import (
	"drpbx/storage-service/file-manager/local"
	"drpbx/storage-service/server"
	"flag"
)

var (
	folder string
	port   string
)

func init() {
	flag.StringVar(&folder, "folder", "./output-folder", "Folder to copy files to")
	flag.StringVar(&port, "port", "8080", "listening port")
}

func main() {
	flag.Parse()

	manager := local.NewLocal(folder)

	svr := server.NewServer(manager, port)

	svr.Serve()

}
