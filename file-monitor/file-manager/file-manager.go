package file_manager

import (
	"drpbx/file-monitor/transfer"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileManager struct {
	folder   string
	transfer transfer.Transfer
	stopped  bool
	oldFiles map[string]os.FileInfo
}

func NewFileManager(folder string, transfer transfer.Transfer) *FileManager {
	return &FileManager{
		folder:   folder,
		transfer: transfer,
		stopped:  false,
		oldFiles: make(map[string]os.FileInfo),
	}
}

func (fm *FileManager) MonitorFolder() {
	for {
		//Eventual consistency is fine. Don't need to poll too often.
		time.Sleep(1 * time.Second)

		if fm.stopped {
			break
		}

		entryMap := map[string]os.FileInfo{}

		entries, err := os.ReadDir(fm.folder)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {

			//Simplifying assumption - ignoring nested folders
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				log.Fatal(err)
			}

			entryMap[entry.Name()] = info

		}

		creates, updates, deletes := fm.diff(entryMap)
		fm.sendDiffs(creates, updates, deletes)

		fm.oldFiles = entryMap

		log.Printf("done")
	}

}

func (fm *FileManager) sendDiffs(creates, updates, deletes []string) {
	for _, create := range creates {
		if err := fm.transfer.Create(filepath.Join(fm.folder, create)); err != nil {
			//This is where one might put retry logic or other error handling
			continue
		}
	}

	for _, update := range updates {
		if err := fm.transfer.Update(filepath.Join(fm.folder, update)); err != nil {
			//This is where one might put retry logic or other error handling
			continue
		}
	}

	for _, dlt := range deletes {
		if err := fm.transfer.Delete(dlt); err != nil {
			//This is where one might put retry logic or other error handling
			continue
		}
	}

}

func (fm *FileManager) diff(entries map[string]os.FileInfo) (create, delete, update []string) {
	for entry, info := range entries {
		oldInfo, ok := fm.oldFiles[entry]

		if !ok {
			create = append(create, entry)
			continue
		}

		if info.ModTime() != oldInfo.ModTime() {
			update = append(update, entry)
			continue
		}
	}

	for entry, _ := range fm.oldFiles {
		if _, ok := entries[entry]; !ok {
			delete = append(delete, entry)
		}
	}

	return create, update, delete
}
