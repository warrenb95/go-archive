package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

var (
	sourceDir   string
	logFilePath string
)

func main() {
	flag.StringVar(&sourceDir, "s", "", "The directory you want to zip up and compress")
	flag.StringVar(&logFilePath, "l", "", "The filepath of the log file, leave empty to lot to stdout")
	flag.Parse()

	if logFilePath != "" {
		logFile, err := os.OpenFile("compressor.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create compressor.log file: %v", err)
		}
		log.SetOutput(logFile)
	}

	if sourceDir == "" {
		log.Fatal("-s source must be specified")
	}

	zipPath := fmt.Sprintf("%s.zip", sourceDir)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Fatalf("creating zip file: %v", err)
	}
	log.Printf("created zip file %s", zipPath)

	w := zip.NewWriter(zipFile)
	mz := myZipper{w: w}
	err = mz.AddFS(os.DirFS(sourceDir))
	if err != nil {
		log.Fatalf("Failed to add FS %s zip writter: %v", sourceDir, err)
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

type myZipper struct {
	w *zip.Writer
}

// AddFS adds the files from fs.FS to the archive.
// It walks the directory tree starting at the root of the filesystem
// adding each file to the zip using deflate while maintaining the directory structure.
func (z *myZipper) AddFS(fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return errors.New("zip: cannot add non-regular file")
		}

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = name
		h.Method = zip.Deflate
		fw, err := z.w.CreateHeader(h)
		if err != nil {
			return err
		}
		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		log.Printf("zipping file: %s, size: %d", info.Name(), info.Size())
		_, err = io.Copy(fw, f)
		return err
	})
}
