package util

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestIsImage(t *testing.T) {
	tests := []struct {
		testName string
		path     string
		expected bool
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
			true,
		},
		{
			"png",
			"../testdata/test2.png",
			true,
		},
		{
			"txt",
			"../testdata/test3.txt",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, _ := os.Open(test.path)
			data := make([]byte, 512, 512)
			f.Read(data)
			actual := isImage(data)
			if actual != test.expected {
				t.Errorf("handler returned wrong error message: got %v want %v", actual, test.expected)
			}
		})
	}
}

func TestGetHash(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			"../testdata/test1.jpg",
			"a2fc620a155852920ef639d81434a7d3",
		},
		{
			"../testdata/test2.png",
			"74bca84a8533a059d5893d6c338daf93",
		},
		{
			"../testdata/test3.txt",
			"68d95371bd37d7466cd78b537cc7a64e",
		},
	}

	for _, test := range tests {
		f, _ := os.Open(test.path)
		actual := getHash(f)
		if actual != test.expected {
			t.Errorf("handler returned wrong error message: got %v want %v", actual, test.expected)
		}
	}
}

func TestCreateImage(t *testing.T) {
	dir, _ := ioutil.TempDir("", "thumbnail_server")
	defer os.RemoveAll(dir)

	tests := []struct {
		testName    string
		path        string
		goldenImage string
		filename    string
		format      string
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
			"../testdata/test1.golden",
			"test1",
			".jpeg",
		},
		{
			"png",
			"../testdata/test2.png",
			"../testdata/test2.golden",
			"test2",
			".png",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, _ := os.Open(test.path)
			actualImagePath := path.Join(dir, test.filename)
			err := createImage(f, actualImagePath)
			if err != nil {
				t.Error(err)
			}
			err = compareImages(test.goldenImage, actualImagePath+test.format)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCreateThumb(t *testing.T) {
	dir, _ := ioutil.TempDir("", "thumbnail_server")
	defer os.RemoveAll(dir)

	tests := []struct {
		testName    string
		path        string
		goldenImage string
		filename    string
		format      string
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
			"../testdata/test1_thumb.golden",
			"thumb_test1",
			".jpeg",
		},
		{
			"png",
			"../testdata/test2.png",
			"../testdata/test2_thumb.golden",
			"thumb_test2",
			".png",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, _ := os.Open(test.path)
			actualImagePath := path.Join(dir, test.filename)
			err := createThumb(f, actualImagePath)
			if err != nil {
				t.Error(err)
			}
			err = compareImages(test.goldenImage, actualImagePath+test.format)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		testName string
		path     string
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
		},
		{
			"png",
			"../testdata/test2.png",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, _ := os.Open(test.path)
			img, format, _ := image.Decode(f)
			buf := new(bytes.Buffer)
			err := encode(buf, format, img)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEncodeAndSave(t *testing.T) {
	dir, _ := ioutil.TempDir("", "thumbnail_server")
	defer os.RemoveAll(dir)

	tests := []struct {
		testName string
		path     string
		filename string
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
			"test1.jpeg",
		},
		{
			"png",
			"../testdata/test2.png",
			"test2.png",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f, _ := os.Open(test.path)
			img, format, _ := image.Decode(f)
			actualImagePath := path.Join(dir, test.filename)
			err := encodeAndSave(img, format, actualImagePath)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestProcessImage(t *testing.T) {
	dir, _ := ioutil.TempDir("", "thumbnail_server")
	defer os.RemoveAll(dir)

	tests := []struct {
		testName    string
		filename    string
		goldenImage string
		goldenThumb string
		image       string
		thumb       string
	}{
		{
			"jpg",
			"../testdata/test1.jpg",
			"../testdata/test1.golden",
			"../testdata/test1_thumb.golden",
			"a2fc620a155852920ef639d81434a7d3.jpeg",
			"thumb_a2fc620a155852920ef639d81434a7d3.jpeg",
		},
		{
			"png",
			"../testdata/test2.png",
			"../testdata/test2.golden",
			"../testdata/test2_thumb.golden",
			"74bca84a8533a059d5893d6c338daf93.png",
			"thumb_74bca84a8533a059d5893d6c338daf93.png",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file, _ := os.Open(test.filename)
			defer file.Close()
			err := ProcessImage(file, dir)
			if err != nil {
				t.Error(err)
			}
			err = compareImages(test.goldenImage, filepath.Join(dir, test.image))
			if err != nil {
				t.Error(err)
			}
			err = compareImages(test.goldenThumb, filepath.Join(dir, test.thumb))
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func compareImages(expectedPath, actualPath string) error {
	actual, err := ioutil.ReadFile(actualPath)
	if err != nil {
		return err
	}
	expected, _ := ioutil.ReadFile(expectedPath)
	if !bytes.Equal(actual, expected) {
		return fmt.Errorf("images are not identical\nactual: %v\nexpected: %v", actual, expected)
	}
	return nil
}
