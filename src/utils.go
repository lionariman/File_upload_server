package src

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

type FileMetadata struct {
	OriginalSize int64 \`json:"original_size"\`
}

func saveMetadata(fileName string, metadata FileMetadata) error {
	targetDir := filepath.Join(DirName, fileName)
	metadataPath := filepath.Join(targetDir, MetadataName)
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return os.WriteFile(metadataPath, data, 0644)
}

func getMetadata(fileName string) (FileMetadata, error) {
	var metadata FileMetadata
	targetDir := filepath.Join(DirName, fileName)
	metadataPath := filepath.Join(targetDir, MetadataName)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return metadata, err
	}
	err = json.Unmarshal(data, &metadata)
	return metadata, err
}

func unboxChunksFromFolder(fileName string) ([][]byte, error) {
	var chunks [][]byte
	targetDir := filepath.Join(DirName, fileName)

	i := 0
	for {
		chunkPath := filepath.Join(targetDir, PortionName+strconv.Itoa(i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			if os.IsNotExist(err) {
				if i == 0 {
					return nil, os.ErrNotExist
				}
				break // End of chunks
			}
			return nil, err
		}
		chunks = append(chunks, data)
		i++
	}
	return chunks, nil
}

func unbzeroChunks(buf [][]byte, originalSize int64) []byte {
	joined := joinChunks(buf)
	if int64(len(joined)) > originalSize {
		return joined[:originalSize]
	}
	return joined
}

func joinChunks(buf [][]byte) []byte {
	var totalLen int
	for _, b := range buf {
		totalLen += len(b)
	}
	joined := make([]byte, 0, totalLen)
	for i := range buf {
		joined = append(joined, buf[i]...)
	}
	return joined
}
