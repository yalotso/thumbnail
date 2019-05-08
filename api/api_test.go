package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/yalotso/thumbnail/config"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestApi(t *testing.T) {
	dir := config.InitTestConfig()
	defer os.RemoveAll(dir)

	t.Run("test multipart", testMultipart)
	t.Run("test base64", testBase64)
	t.Run("test reference", testReference)
}

func testMultipart(t *testing.T) {
	url := "/multipart"
	r := routing.New()
	r.Post(url, Multipart)

	t.Run("test image file", func(t *testing.T) {
		buf := new(bytes.Buffer)
		fw := multipart.NewWriter(buf)

		formFile, _ := fw.CreateFormFile("image", "test1.jpg")
		file, _ := os.Open("../testdata/test1.jpg")
		io.Copy(formFile, file)
		file.Close()
		fw.Close()

		req, _ := http.NewRequest("POST", url, buf)
		req.Header.Set("Content-Type", fw.FormDataContentType())

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		if status := res.Code; status != http.StatusOK {
			t.Error(res.Body.String())
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test non-image file", func(t *testing.T) {
		buf := new(bytes.Buffer)
		fw := multipart.NewWriter(buf)

		formFile, _ := fw.CreateFormFile("image", "test3.txt")
		file, _ := os.Open("../testdata/test3.txt")
		io.Copy(formFile, file)
		file.Close()
		fw.Close()

		req, _ := http.NewRequest("POST", url, buf)
		req.Header.Set("Content-Type", fw.FormDataContentType())

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "test3.txt: unknown format"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})

	t.Run("test empty body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "multipart/form-data")

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "missing form body"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})
}

func testBase64(t *testing.T) {
	url := "/base64"
	r := routing.New()
	r.Post(url, Base64)

	t.Run("test base64 image", func(t *testing.T) {
		file, _ := os.Open("../testdata/test3.txt")
		data, _ := ioutil.ReadAll(file)
		file.Close()

		jsonData, _ := json.Marshal(string(data))

		req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test base64 non-image", func(t *testing.T) {
		file, _ := os.Open("../testdata/test4.txt")
		data, _ := ioutil.ReadAll(file)
		file.Close()

		jsonData, _ := json.Marshal(string(data))

		req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "unknown format"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})

	t.Run("test empty body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "missing request body"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})
}

func testReference(t *testing.T) {
	url := "/reference"
	r := routing.New()
	r.Get(url, Reference)

	t.Run("test image url", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)

		q := req.URL.Query()
		q.Add("url", "https://cdn-images-1.medium.com/max/2400/1*30aoNxlSnaYrLhBT0O1lzw.png")
		req.URL.RawQuery = q.Encode()

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test non-image url", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)

		q := req.URL.Query()
		q.Add("url", "https://github.com/golang/go/blob/master/api/go1.12.txt")
		req.URL.RawQuery = q.Encode()

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "unknown format"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})

	t.Run("test empty body", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)

		q := req.URL.Query()
		q.Add("url", "")
		req.URL.RawQuery = q.Encode()

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		expected := "missing url parameter"
		if body := res.Body.String(); body[:len(body)-1] != expected {
			t.Errorf("handler returned wrong error message: got %v want %v", body, expected)
		}
	})
}
