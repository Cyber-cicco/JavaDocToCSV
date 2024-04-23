package main

import (
	"io"
	"log"
	"net/http"

	"github.com/Cyber-cicco/javaToCSV/decoding"
)

const BASE_LANG = "https://docs.oracle.com/javase/8/docs/api/java/lang/"

const JAVALANG_URL = "https://docs.oracle.com/javase/8/docs/api/java/lang/package-summary.html"

func main() {
	resp, err := http.Get(JAVALANG_URL)

	if err != nil {
		log.Fatalln(err)
	}

    content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

    decoding.FindLinksInSummary(content)
}

