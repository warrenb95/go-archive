package main

import (
	"log"
	"os"
)

func main() {
	sourceFile, err := os.Open("test.txt")
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Fatalf("failed to close source file: %v", err)
		}
	}()
	if err != nil {
		log.Fatalf("failed to open test jpeg: %v", err)
	}
}
