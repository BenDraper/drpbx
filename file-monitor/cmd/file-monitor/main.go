package main

import (
	file_manager "drpbx/file-monitor/file-manager"
	http_transfer "drpbx/file-monitor/transfer/http"
	"flag"
)

var (
	folder        string
	writerAddress string
)

func init() {
	flag.StringVar(&folder, "folder", "./input-folder", "Folder to monitor")
	flag.StringVar(&writerAddress, "address", "http://localhost:8080/", "Address of the file writer")
}

func main() {
	flag.Parse()

	transfer := http_transfer.NewHTTPTransfer(writerAddress)

	manager := file_manager.NewFileManager(folder, transfer)

	manager.MonitorFolder()

}
