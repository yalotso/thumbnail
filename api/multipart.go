package api

import (
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	. "github.com/yalotso/thumbnail/config"
	"github.com/yalotso/thumbnail/util"
	"mime/multipart"
)

func Multipart(c *routing.Context) error {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	for _, fileHeaders := range c.Request.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			err = uploadImage(fileHeader)
			if err != nil {
				return fmt.Errorf("%s: %s", fileHeader.Filename, err)
			}
		}
	}
	return nil
}

func uploadImage(fh *multipart.FileHeader) error {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	err = util.ProcessImage(file, Conf.MultipartDir)
	if err != nil {
		return err
	}
	return nil
}
