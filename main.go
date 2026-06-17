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
	mux := http.NewServeMux()
	mux.HandleFunc("/", processor)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    src.PortName,
		Handler: mux,
	}

	go func() {
		log.Println("Server started on", src.PortName)
		fmt.Println(src.Endpoints)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: [%s]\n", err)
		}
	}()

	<-done
	log.Println("\nServer stopping")
}

// Router
func processor(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/upload_file" {
		src.UploadFile(w, r)
	} else if strings.HasPrefix(path, "/delete_file:") {
		fileName := strings.TrimPrefix(path, "/delete_file:")
		src.DeleteFile(w, r, fileName)
	} else if strings.HasPrefix(path, "/get_file:") {
		fileName := strings.TrimPrefix(path, "/get_file:")
		src.GetFile(w, r, fileName)
	} else if path == "/delete_all_files" {
		src.DeleteAllFiles(w, r)
	} else if path == "/shutdown" {
		src.ServerShutdown(w)
	} else {
		fmt.Println("Wrong endpoint! ->", path)
		http.NotFound(w, r)
	}
}
