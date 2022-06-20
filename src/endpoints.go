package src

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Endpoint: /delete_file:fileName.
// Функция DeleteFile удаляет папку с названием файла,
// в которой лежат ее мегабайты, с сервера
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		fmt.Printf("Wrong request method -> [%s]\n", r.Method)
		return
	}
	if strings.ContainsAny(Endpoint, ":") == false {
		fmt.Printf("Empty item after [%s:___?___]\n", Endpoint)
		return
	}
	var currFileName string = strings.Split(Endpoint, ":")[1]
	fmt.Printf("[%s] delete [%s]\n", r.Method, currFileName)
	PortionDirName += currFileName
	var template string = DirName + "/" + PortionDirName
	err := os.RemoveAll(template)
	if err != nil {
		fmt.Printf("Cannot remove [%s] directory\n", template)
		log.Println(err)
		return
	}
	PortionDirName = PortionDirNameTmp
}

// Endpoint: /get_file.
// Функция GetFile восттанавливает файл из частей (мегабайтов) и возвращает его
func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Printf("Wrong request method -> [%s]\n", r.Method)
		return
	}
	if strings.ContainsAny(Endpoint, ":") == false {
		fmt.Printf("Empty item after [%s:___?___]\n", Endpoint)
		return
	}
	var currFileName string = strings.Split(Endpoint, ":")[1]
	fmt.Printf("[%s] get [%s]\n", r.Method, currFileName)
	PortionDirName += currFileName
	unboxedChunks, err := unboxChunksFromFolder()
	if err != nil {
		fmt.Printf("Cannot unbox chunks from [%s]\n", PortionDirName)
		log.Println(err)
		return
	}
	unbzeroChunks(unboxedChunks)
	joinedChunks := joinChunks(unboxedChunks)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(joinedChunks)
}

// Endpoint: /delete_all_files.
// Функция DeleteAllFiles удаляет все папки с мегабайтами файлов с сервера
func DeleteAllFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		fmt.Printf("Wrong request method -> [%s]\n", r.Method)
		return
	}
	fmt.Printf("[%s] delete all files in [%s]\n", r.Method, DirName)
	err := os.RemoveAll(DirName + "/")
	if err != nil {
		fmt.Printf("cannot remove [%s]\n", DirName)
		log.Println(err)
		return
	}
}

// Endpoint: /upload_file.
// Функция UploadFile принимает файл и делит его по частям на мегабайты
// и сохраняет эти части в отдельной папке с названием файла
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Printf("Wrong request method -> [%s]\n", r.Method)
		return
	}
	fmt.Printf("[%s] upload [%s]\n", r.Method, r.URL.String())
	r.ParseMultipartForm(500 << 20) // 500 mb
	var buf bytes.Buffer
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	PortionDirName += strings.Split(fileHeader.Filename, ".")[0]
	fmt.Printf("FOLDER [%s]\n", DirName+"/"+PortionDirName)
	numberOfBytesCopied, err := io.Copy(&buf, file)
	if err != nil {
		fmt.Printf("%s file has %d bytes (too big file)\n",
			fileHeader.Filename, numberOfBytesCopied)
		return
	}
	chunkedBuf := makeChunks(buf.Bytes())
	bzeroChunks(chunkedBuf)
	saveChunksIntoFolder(chunkedBuf)
}

// Endpoint: /shutdown.
// Функция ServerShutdown выключает сервер
func ServerShutdown(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	os.Exit(0)
}
