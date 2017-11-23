package main

import (
	"bufio"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/urfave/cli"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func resizeJpeg(inName, outName string, size int) error {
	file, err := os.Open(inName)
	defer file.Close()
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	var m image.Image
	if img.Bounds().Size().X > img.Bounds().Size().Y {
		m = resize.Resize(uint(size), 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, uint(size), img, resize.Lanczos3)
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

func resizeOperation(origpath, resizepath string, size int) error {
	//printOperation(origpath, resizepath)
	err := createPathToFile(resizepath)
	if err != nil {
		return err
	}
	err = resizeJpeg(origpath, resizepath, size)
	if err != nil {
		return err
	}
	return nil
}

var operationCount int = 0

func NewVisitFunc(operation func(string, string, int) error, origRoot, resizeRoot string, size int) func(string, os.FileInfo, error) error {

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
			err := operation(path, resizepath, size)
			if err != nil {
				printOperation(path, resizepath)
				log.Printf("Error: %q", err.Error())
				return nil // Skip error!
			}
			operationCount++
			fmt.Printf("Converted: %d images.\r", operationCount)
		}
		return nil
	}
}

func resizeTree(origRoot, resizeRoot string, size int) {
	//visit := NewVisitFunc(printOperation, origRoot, resizeRoot)
	visit := NewVisitFunc(resizeOperation, origRoot, resizeRoot, size)
	//err := filepath.Walk(origRoot, visit)
	filepath.Walk(origRoot, visit)
	//log.Printf("Visited %d images, Walk() returned: %v", operationCount, err)
}
func askUserToContinue() bool {
	fmt.Printf("Continue? [yN]: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil || scanner.Text() != "y" {
		return false
	}
	return true
}

func mainCommand(c *cli.Context) error {
	source := c.Args().Get(0)
	dest := c.String("dest")
	size := c.Int("size")
	if source == "" {
		log.Fatal("Source must be given.")
	}
	if _, err := os.Stat(source); err != nil {
		log.Fatal(err)
	}

	log.Printf("Converting: %q to %q with size %d. \n", source, dest, size)

	if !c.Bool("no-ask") && !askUserToContinue() {
		log.Printf("Cancelled by user.")
		return nil
	}
	resizeTree(source, dest, size)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "resizr"
	app.Usage = "Create small image previews in seperate folder structure\n\n   Example: resizr --dest /home/user/preview --size 1024 /home/user/pictures"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dest, o",
			Value: ".",
			Usage: "Destination directory",
		},
		cli.IntFlag{
			Name:  "size, s",
			Value: 1024,
			Usage: "Set default max image width/height",
		},
		cli.BoolFlag{
			Name:  "no-ask, y",
			Usage: "Skip question if we should continue.",
		},
	}
	app.Action = mainCommand
	app.Run(os.Args)
}
