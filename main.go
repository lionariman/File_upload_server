package main

import (
	"fileChunker/src"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	handler := http.HandlerFunc(processor)
	http.Handle("/", handler)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := http.ListenAndServe(src.PortName, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: [%s]\n", err)
		}
	}()
	log.Println("Server started")
	fmt.Println(src.Endpoints)
	<-done
	log.Println("\nServer stopped")
	defer close(done)
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
	} else if src.Endpoint == "/shutdown" {
		src.ServerShutdown(w)
	} else {
		fmt.Println("Wrong endpoint! ->", src.Endpoint)
	}
}
