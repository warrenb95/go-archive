package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.Create("compressor.log")
	if err != nil {
		log.Fatalf("failed to create compressor.log file: %v", err)
	}
	log.SetOutput(logFile)

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
