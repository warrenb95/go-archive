package compresser

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func Compress(outName string, source io.Reader) error {
	destFile, err := os.Create(outName)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return fmt.Errorf("creating destination: %w", err)
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			log.Printf("failed to close destFile: %v", err)
		}
	}()

	gw := gzip.NewWriter(destFile)
	defer func() {
		if err := gw.Close(); err != nil {
			log.Printf("failed to close gzip writer: %v", err)
		}
	}()

	b, err := io.Copy(gw, source)
	if err != nil {
		log.Printf("failed to copy and compress source: %v ", err)
		return fmt.Errorf("copying source to destination: %w", err)
	}
	log.Printf("bytes copied: %d ", b)

	if err := gw.Flush(); err != nil {
		log.Printf("flushing gzip writer: %v", err)
		return fmt.Errorf("flushing gzip writer: %w", err)
	}

	return nil
}
