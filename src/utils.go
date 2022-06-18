package src

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

func makeChunks(buf []byte) [][]byte {
	var chunkedBuf [][]byte
	var first, last int
	for i := 0; i < len(buf)/ChunkSize+1; i++ {
		first = i * ChunkSize
		last = i*ChunkSize + ChunkSize
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
	for ChunkSize > len(buf[lastElem]) {
		buf[lastElem] = append(buf[lastElem], 0)
	}
}

func createDirectories() {
	err := os.MkdirAll(DirName+"/"+PortionDirName, 0777)
	if os.IsExist(err) == true {
		fmt.Printf("[%s] already exists\n", DirName)
	} else if err != nil {
		log.Fatal(err)
	}
	PortionDirName = PortionDirNameTmp
}

func saveChunksIntoFolder(buf [][]byte) {
	var template string = DirName + "/" + PortionDirName + "/" + PortionName
	createDirectories()
	for i := range buf {
		err := os.WriteFile(template+strconv.Itoa(i), buf[i], 0777)
		if err != nil {
			log.Fatal()
		}
	}
}

func unboxChunksFromFolder() ([][]byte, error) {
	var chunks [][]byte
	var template string = DirName + "/" + PortionDirName + "/" + PortionName
	var i int = 0
	for {
		data, err := os.ReadFile(template + strconv.Itoa(i))
		if os.IsExist(err) == false {
			fmt.Printf("[%s] Directory is not exist\n", template)
			return nil, err
		} else if err != nil {
			return nil, err
		}
		chunks = append(chunks, data)
		i++
	}
	// PortionDirName = PortionDirNameTmp
	// return chunks, nil
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
	err := os.WriteFile("new_"+FileName, jbuf, 0666)
	if err != nil {
		log.Fatal()
	}
}
