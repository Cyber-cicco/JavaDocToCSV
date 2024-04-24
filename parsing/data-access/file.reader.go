package dataaccess

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

func parseFolders(files []fs.DirEntry, path, suffix string, executable func(fileContent []byte, filePath string)) {
	for _, file := range files {
		if file.IsDir() {
			files, err := os.ReadDir(path + "/" + file.Name())
			if err != nil {
				log.Fatalf("Got error %s", err)
			}
			parseFolders(files, path+"/"+file.Name(), suffix, executable)
		} else if strings.HasSuffix(file.Name(), suffix) {
			filePath := path + "/" + file.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Got error %s", err)
			}
			executable(content, filePath)
		}
	}
}

func ParseFolders(suffix, path string, executable func(fileContent []byte, filePath string)) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Got error %s", err)
	}
	parseFolders(files, path, suffix, executable)
}
