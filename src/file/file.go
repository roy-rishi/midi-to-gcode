package file

import (
	"log"
	"os"
)

func ReadBin(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}
