package main

import (
	"log"
	"net/http"
)

func serveAssets() {
	mux := http.NewServeMux()

	fileSystemHandler := http.FileServer(http.Dir("assets/"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fileSystemHandler))

	log.Fatal(http.ListenAndServe(":3000", mux))
}
