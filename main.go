package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	chunkSize   int    = 1024 * 1024 // 1 mb
	fileName    string = "file.txt"
	dirName     string = "1mb_portions"
	portionName string = "1mb_portion_"
	portName    string = ":8080"
)

func main() {
	fmt.Println("Server is listening...\n" +
		"/upload -> upload file to server\n" +
		"/file -> get file from server")
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/file", getFile)
	err := http.ListenAndServe(portName, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	unboxedChunks := unboxChunksFromFolder()
	unbzeroChunks(unboxedChunks)
	joinedChunks := joinChunks(unboxedChunks)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/octet-stream")
	w.Write(joinedChunks)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	r.ParseMultipartForm(500 << 20) // 500 mb
	var buf bytes.Buffer
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	numberOfBytesCopied, err := io.Copy(&buf, file)
	if err != nil {
		fmt.Printf("%s file has %d bytes (too big file)\n", fileHeader.Filename, numberOfBytesCopied)
		return
	}
	chunkedBuf := makeChunks(buf.Bytes())
	bzeroChunks(chunkedBuf)
	saveChunksIntoFolder(chunkedBuf)
}

//--------------
func makeChunks(buf []byte) [][]byte {
	var chunkedBuf [][]byte
	var first, last int
	for i := 0; i < len(buf)/chunkSize+1; i++ {
		first = i * chunkSize
		last = i*chunkSize + chunkSize
		if last > len(buf) {
			last = len(buf)
		}
		chunkedBuf = append(chunkedBuf, buf[first:last])
	}
	return chunkedBuf
}

func unbzeroChunks(buf [][]byte) {
	var lastElem int = len(buf) - 1
	for i := 0; i <= len(buf[lastElem]); i++ {
		if buf[lastElem][i] == 0x00 {
			buf[lastElem] = buf[lastElem][:i]
		}
	}
}

func bzeroChunks(buf [][]byte) {
	for i := range buf {
		for chunkSize > len(buf[i]) {
			buf[i] = append(buf[i], 0)
		}
	}
}

func saveChunksIntoFolder(buf [][]byte) {
	var template string = dirName + "/" + portionName
	err := os.Mkdir(dirName, 0777)
	if os.IsExist(err) == true {
		fmt.Println("[", dirName, "] already exists")
	} else if err != nil {
		log.Fatal(err)
	}
	for i := range buf {
		err := os.WriteFile(template+strconv.Itoa(i), buf[i], 0777)
		if err != nil {
			log.Fatal()
		}
	}
}

func unboxChunksFromFolder() [][]byte {
	var chunks [][]byte
	var template string = dirName + "/" + portionName
	var i int = 0
	for {
		data, err := os.ReadFile(template + strconv.Itoa(i))
		if err != nil {
			break
		}
		chunks = append(chunks, data)
		i++
	}
	return chunks
}

func joinChunks(buf [][]byte) []byte {
	var joined []byte
	for i := range buf {
		for j := range buf[i] {
			joined = append(joined, buf[i][j])
		}
	}
	return joined
}

func createFileFromChunks(buf [][]byte) {
	jbuf := joinChunks(buf)
	err := os.WriteFile("new_"+fileName, jbuf, 0666)
	if err != nil {
		log.Fatal()
	}
}

//--------------
