package main

import (
	"github.com/Cyber-cicco/HTMLtoDB/config"
	dataaccess "github.com/Cyber-cicco/HTMLtoDB/data-access"
)



func main() {

    dataaccess.ParseFolders(".html", config.URL_RESOURCES, func(fileContent, filePath string) {

    })
}

