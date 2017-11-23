package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func resizeJpeg(inName, outName string) error {
	file, err := os.Open(inName)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	var m image.Image
	if img.Bounds().Size().X > img.Bounds().Size().Y {
		m = resize.Resize(1024, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 1024, img, resize.Lanczos3)
	}

	out, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)

	return nil
}

func createPathToFile(resizepath string) error {
	folderPath := filepath.Dir(resizepath)
	return os.MkdirAll(folderPath, os.ModePerm) // Returns an error (or nil)
}

func printOperation(origpath, resizepath string) {
	log.Printf("Operation: %s -> %s \n", origpath, resizepath)
}

func resizeOperation(origpath, resizepath string) error {
	printOperation(origpath, resizepath)
	err := createPathToFile(resizepath)
	if err != nil {
		return err
	}
	err = resizeJpeg(origpath, resizepath)
	if err != nil {
		return err
	}
	return nil
}

var operationCount int = 0

func NewVisitFunc(operation func(string, string) error, origRoot, resizeRoot string) func(string, os.FileInfo, error) error {

	return func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if err != nil { // passed from walk
			return err
		}
		ext := strings.ToLower(string(filepath.Ext(path)))
		if ext == ".jpeg" || ext == ".jpg" {
			relativPath, _ := filepath.Rel(origRoot, path)
			resizepath := filepath.Join(resizeRoot, relativPath)
			err := operation(path, resizepath)
			if err != nil {
				return err
			}
			operationCount++
		}
		return nil
	}
}

func resizeTree(origRoot, resizeRoot string) {
	//visit := NewVisitFunc(printOperation, origRoot, resizeRoot)
	visit := NewVisitFunc(resizeOperation, origRoot, resizeRoot)
	err := filepath.Walk(origRoot, visit)
	log.Printf("Visited %d images, Walk() returned: %v", operationCount, err)
}

func main() {
	resizeTree("/home/markus/test/testbilder_orig", "test/resized/brenta")
}
