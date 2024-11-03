package main

import (
	"archive/zip"
	"log"
	"os"

	"github.com/warrenb95/go-archive/internal/compresser"
)

func main() {
	logFile, err := os.OpenFile("compressor.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create compressor.log file: %v", err)
	}
	log.SetOutput(logFile)

	zipFile, err := os.Create("archive.zip")
	if err != nil {
		log.Fatalf("creating zip file: %v", err)
	}

	w := zip.NewWriter(zipFile)
	err = w.AddFS(os.DirFS("./internal/compresser/test_data"))
	if err != nil {
		log.Fatalf("Failed to add e/archive to zip writter: %v", err)
	}
	err = w.Close()
	if err != nil {
		log.Fatalf("closing zip writer: %v", err)
	}

	err = zipFile.Close()
	if err != nil {
		log.Fatalf("failed to close zip file: %v", err)
	}

	zipFile, err = os.Open("archive.zip")
	if err != nil {
		log.Fatalf("opening zip file: %v", err)
	}

	err = compresser.Compress("archive.zip.gzip", zipFile)
	if err != nil {
		log.Fatalf("compressing archive file: %v", err)
	}

	err = zipFile.Close()
	if err != nil {
		log.Fatalf("failed to close zip file: %v", err)
	}
}
