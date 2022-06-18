package main

import (
	"fileChunker/src"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	handler := http.HandlerFunc(processor)
	http.Handle("/", handler)
	fmt.Println(src.Endpoints)
	log.Println("listening on :8080...")
	err := http.ListenAndServe(src.PortName, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func processor(w http.ResponseWriter, r *http.Request) {
	src.Endpoint = r.URL.String()
	if src.Endpoint == "/upload_file" {
		src.UploadFile(w, r)
	} else if strings.Split(src.Endpoint, ":")[0] == "/delete_file" {
		src.DeleteFile(w, r)
	} else if strings.Split(src.Endpoint, ":")[0] == "/get_file" {
		src.GetFile(w, r)
	} else if src.Endpoint == "/delete_all_files" {
		src.DeleteAllFiles(w, r)
	} else {
		fmt.Println("Wrong endpoint! ->", src.Endpoint)
	}
}
