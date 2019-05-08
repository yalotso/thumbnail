package util

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func ProcessImage(reader io.Reader, dir string) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if !isImage(data) {
		return errors.New("unknown format")
	}
	rs := bytes.NewReader(data)
	hash := getHash(rs)
	err = createImage(rs, path.Join(dir, hash))
	if err != nil {
		return err
	}
	err = createThumb(rs, path.Join(dir, "thumb_"+hash))
	if err != nil {
		return err
	}
	return nil
}

func isImage(data []byte) bool {
	if len(data) < 512 {
		return false
	}
	mimeType := http.DetectContentType(data[:512])
	return mimeType == "image/png" || mimeType == "image/jpeg"
}

func getHash(reader io.ReadSeeker) string {
	h := md5.New()
	io.Copy(h, reader)
	reader.Seek(0, 0)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createImage(reader io.ReadSeeker, filename string) error {
	img, format, err := image.Decode(reader)
	if err != nil {
		return err
	}
	defer reader.Seek(0, 0)
	filename += "." + format
	err = encodeAndSave(img, format, filename)
	if err != nil {
		return err
	}
	return nil
}

func createThumb(reader io.ReadSeeker, filename string) error {
	img, format, err := image.Decode(reader)
	if err != nil {
		return err
	}
	defer reader.Seek(0, 0)
	thumb := resize.Thumbnail(100, 100, img, resize.NearestNeighbor)
	filename += "." + format
	err = encodeAndSave(thumb, format, filename)
	if err != nil {
		return err
	}
	return nil
}

func encodeAndSave(img image.Image, format, path string) error {
	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()
	err := encode(buf, format, img)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func encode(dst io.Writer, format string, img image.Image) error {
	switch format {
	case "jpeg":
		return jpeg.Encode(dst, img, nil)
	case "png":
		return png.Encode(dst, img)
	}
	return errors.New("unknown format")
}
