package main

import (
	"io"
	"log"
	"net/http"
	"os"

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

	linkMap := decoding.FindLinksInSummary(content)

	for classType, links := range linkMap {
		for _, link := range links {

			err := os.MkdirAll("./resources/lang/"+classType, 0777)

			if err != nil {
				log.Fatalln(err)
			}

            resp, err = http.Get(BASE_LANG + "/" + link)
			if err != nil {
				log.Fatalln(err)
			}

            body, err := io.ReadAll(resp.Body)

			if err != nil {
				log.Fatalln(err)
			}

            if err = os.WriteFile("./resources/lang/" + classType + "/" + link, body, 0644); err != nil {
				log.Fatalln(err)
            }
		}
	}
}
