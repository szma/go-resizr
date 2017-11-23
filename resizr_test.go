package main

import "testing"
import "os"
import "fmt"

func TestResize(t *testing.T) {
	origfile := "test/orig/test.jpg"
	resizefile := "test/resized/test_resized.jpg"

	// Delete last test result and check that this was successful
	err := os.Remove(resizefile)
	if err != nil {
		//t.Errorf(err.Error())
		fmt.Println(err.Error())
	}
	if _, err := os.Stat(resizefile); !os.IsNotExist(err) {
		t.Fail()
	}
	createPathToFile(resizefile)
	// Do the resizing
	resizeJpeg(origfile, resizefile, 128)
	// Check if a new file is there
	if _, err := os.Stat(resizefile); os.IsNotExist(err) {
		t.Fail()
	}
	// Check if a original file is there
	if _, err := os.Stat(origfile); os.IsNotExist(err) {
		t.Fail()
	}
	stat_orig, _ := os.Stat(origfile)
	stat_resize, _ := os.Stat(resizefile)
	if stat_orig.Size() == stat_resize.Size() {
		t.Fail()
	}
	if stat_orig.Size() == 0 || stat_resize.Size() == 0 {
		t.Fail()
	}

}

func TestResizeTree(t *testing.T) {
	ResizeTree("test/orig", "test/resized", 128)
	//resizeTree("/home/markus/test/2013-06-09 Brenta", "test/resized/brenta")
}
