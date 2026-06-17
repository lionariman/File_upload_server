package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// DeleteFile deletes the directory associated with the file name.
func DeleteFile(w http.ResponseWriter, r *http.Request, fileName string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName = filepath.Base(fileName)
	if fileName == "" || fileName == "." || fileName == "/" {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	fmt.Printf("[%s] delete [%s]\n", r.Method, fileName)
	targetDir := filepath.Join(DirName, fileName)
	err := os.RemoveAll(targetDir)
	if err != nil {
		fmt.Printf("Cannot remove [%s] directory: %v\n", targetDir, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetFile restores the file from chunks and returns it.
func GetFile(w http.ResponseWriter, r *http.Request, fileName string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName = filepath.Base(fileName)
	if fileName == "" || fileName == "." || fileName == "/" {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	fmt.Printf("[%s] get [%s]\n", r.Method, fileName)
	metadata, err := getMetadata(fileName)
	if err != nil {
		fmt.Printf("Cannot get metadata for [%s]: %v\n", fileName, err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	targetDir := filepath.Join(DirName, fileName)
	// Check if all chunks exist before starting stream
	for i := 0; i < int((metadata.OriginalSize+int64(ChunkSize)-1)/int64(ChunkSize)); i++ {
		chunkPath := filepath.Join(targetDir, PortionName+strconv.Itoa(i))
		if _, err := os.Stat(chunkPath); err != nil {
			http.Error(w, "Missing file chunks", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	var bytesRemaining int64 = metadata.OriginalSize
	chunkIndex := 0
	for bytesRemaining > 0 {
		chunkPath := filepath.Join(targetDir, PortionName+strconv.Itoa(chunkIndex))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			fmt.Printf("Error opening chunk %d for %s: %v\n", chunkIndex, fileName, err)
			break
		}

		limit := int64(ChunkSize)
		if bytesRemaining < limit {
			limit = bytesRemaining
		}

		n, err := io.CopyN(w, chunkFile, limit)
		bytesRemaining -= n
		chunkFile.Close()

		if err != nil && err != io.EOF {
			fmt.Printf("Error streaming chunk %d for %s: %v\n", chunkIndex, fileName, err)
			break
		}
		chunkIndex++
	}
}

// DeleteAllFiles deletes all uploaded files.
func DeleteAllFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("[%s] delete all files in [%s]\n", r.Method, DirName)
	err := os.RemoveAll(DirName)
	if err != nil {
		fmt.Printf("cannot remove [%s]: %v\n", DirName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UploadFile receives a file, chunks it, and saves it.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("[%s] upload [%s]\n", r.Method, fileHeader.Filename)
	fileName := filepath.Base(fileHeader.Filename)
	targetDir := filepath.Join(DirName, fileName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var originalSize int64
	chunkIndex := 0
	for {
		chunkPath := filepath.Join(targetDir, PortionName+strconv.Itoa(chunkIndex))
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		n, err := io.CopyN(chunkFile, file, int64(ChunkSize))
		originalSize += n

		if err == io.EOF {
			// Padding the last chunk
			if n < int64(ChunkSize) {
				padding := make([]byte, int64(ChunkSize)-n)
				chunkFile.Write(padding)
			}
			chunkFile.Close()
			break
		} else if err != nil {
			chunkFile.Close()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		chunkFile.Close()
		chunkIndex++
	}

	err = saveMetadata(fileName, FileMetadata{OriginalSize: originalSize})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ServerShutdown shuts down the server.
func ServerShutdown(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("Shutting down...")
	os.Exit(0)
}
