package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	chunkSize         int    = 1 << 20 // 1 mb
	fileName          string = "file.txt"
	dirName           string = "portions"
	portionDirName    string = "1mb_"
	portionDirNameTmp string = "1mb_"
	portionName       string = "1mb_"
	portName          string = ":8080"
	endpoint          string = ""
	endpoints         string = "\nENDPOINTS\n\n" +
		"/upload_file\n" +
		"get_file:{file name}\n" +
		"delete_file:{file name}\n" +
		"delete_all_files\n"
)

func main() {
	handler := http.HandlerFunc(processor)
	http.Handle("/", handler)
	fmt.Println(endpoints)
	log.Println("listening on :8080...")
	err := http.ListenAndServe(portName, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func processor(w http.ResponseWriter, r *http.Request) {
	endpoint = r.URL.String()
	if endpoint == "/upload_file" {
		uploadFile(w, r)
	} else if strings.Split(endpoint, ":")[0] == "/delete_file" {
		deleteFile(w, r)
	} else if strings.Split(endpoint, ":")[0] == "/get_file" {
		getFile(w, r)
	} else if endpoint == "/delete_all_files" {
		deleteAllFiles(w, r)
	} else {
		fmt.Println("Wrong endpoint! ->", endpoint)
	}
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	var currFileName string = strings.Split(endpoint, ":")[1]
	fmt.Printf("[%s] delete [%s]\n", r.Method, currFileName)
	portionDirName += currFileName
	var template string = dirName + "/" + portionDirName
	err := os.RemoveAll(template)
	if err != nil {
		fmt.Printf("Cannot remove [%s] directory\n", template)
		log.Fatal(err)
	}
	portionDirName = portionDirNameTmp
}

func getFile(w http.ResponseWriter, r *http.Request) {
	var currFileName string = strings.Split(endpoint, ":")[1]
	fmt.Printf("[%s] get [%s]\n", r.Method, currFileName)
	portionDirName += currFileName
	unboxedChunks := unboxChunksFromFolder()
	unbzeroChunks(unboxedChunks)
	joinedChunks := joinChunks(unboxedChunks)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(joinedChunks)
}

func deleteAllFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] delete all files in [%s]\n", r.Method, dirName)
	err := os.RemoveAll(dirName + "/")
	if err != nil {
		fmt.Printf("cannot remove [%s]\n", dirName)
		log.Fatal(err)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] upload [%s]\n", r.Method, r.URL.String())
	r.ParseMultipartForm(500 << 20) // 500 mb
	var buf bytes.Buffer
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	portionDirName += strings.Split(fileHeader.Filename, ".")[0]
	fmt.Printf("FOLDER [%s]\n", dirName+"/"+portionDirName)
	numberOfBytesCopied, err := io.Copy(&buf, file)
	if err != nil {
		fmt.Printf("%s file has %d bytes (too big file)\n", fileHeader.Filename, numberOfBytesCopied)
		return
	}
	chunkedBuf := makeChunks(buf.Bytes())
	bzeroChunks(chunkedBuf)
	saveChunksIntoFolder(chunkedBuf)
}

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
	buf[lastElem] = bytes.Trim(buf[lastElem], "\x00")
}

func bzeroChunks(buf [][]byte) {
	var lastElem int = len(buf) - 1
	for chunkSize > len(buf[lastElem]) {
		buf[lastElem] = append(buf[lastElem], 0)
	}
}

func createDirectories() {
	err := os.MkdirAll(dirName+"/"+portionDirName, 0777)
	if os.IsExist(err) == true {
		fmt.Printf("[%s] already exists\n", dirName)
	} else if err != nil {
		log.Fatal(err)
	}
	portionDirName = portionDirNameTmp
}

func saveChunksIntoFolder(buf [][]byte) {
	var template string = dirName + "/" + portionDirName + "/" + portionName
	createDirectories()
	for i := range buf {
		err := os.WriteFile(template+strconv.Itoa(i), buf[i], 0777)
		if err != nil {
			log.Fatal()
		}
	}
}

func unboxChunksFromFolder() [][]byte {
	var chunks [][]byte
	var template string = dirName + "/" + portionDirName + "/" + portionName
	var i int = 0
	for {
		data, err := os.ReadFile(template + strconv.Itoa(i))
		if err != nil {
			break
		}
		chunks = append(chunks, data)
		i++
	}
	portionDirName = portionDirNameTmp
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
