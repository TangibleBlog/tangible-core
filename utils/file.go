package utils

import (
	"log"
	"os"
)

func OpenFile(filename string) []byte {
	result, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func UnCheckOpenFile(filename string) []byte {
	result, err := os.ReadFile(filename)
	if err != nil {
		return []byte("")
	}
	return result
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
