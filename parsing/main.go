package main

import (
	"github.com/Cyber-cicco/HTMLtoDB/config"
	dataaccess "github.com/Cyber-cicco/HTMLtoDB/data-access"
	"github.com/Cyber-cicco/HTMLtoDB/decoding"
)



func main() {
    dataaccess.ParseFolders(".html", config.URL_RESOURCES, decoding.ParseSingleFile)
}

