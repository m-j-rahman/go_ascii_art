package main

import (
	"asciiweb/asciiweb"
	"fmt"
	"log"
	"net/http"
)

const Port string = "3000"

func main() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.HandleFunc("/", asciiweb.MainPage)
	http.HandleFunc("/ascii-art", asciiweb.AsciiArt)
	http.HandleFunc("/export", asciiweb.ExportFile)

	fmt.Printf("Starting server...\n")
	fmt.Printf("Listening on port 3000!\n")

	asciiweb.OpenBrowser("http://localhost:3000")

	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatal(err)
	}
}
