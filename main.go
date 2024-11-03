package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"
	"os"
)

var sourceDir string

func main() {
	flag.StringVar(&sourceDir, "s", "", "The directory you want to zip up and compress")
	flag.Parse()

	logFile, err := os.OpenFile("compressor.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create compressor.log file: %v", err)
	}
	log.SetOutput(logFile)

	if sourceDir == "" {
		log.Fatal("-s source and -o output need to be specified")
	}

	zipPath := fmt.Sprintf("%s.zip", sourceDir)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Fatalf("creating zip file: %v", err)
	}
	log.Printf("created zip file %s", zipPath)

	w := zip.NewWriter(zipFile)
	err = w.AddFS(os.DirFS(sourceDir))
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
}
