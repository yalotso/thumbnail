package api

import (
	"errors"
	"github.com/go-ozzo/ozzo-routing"
	. "github.com/yalotso/thumbnail/config"
	"github.com/yalotso/thumbnail/util"
	"net/http"
)

func Reference(c *routing.Context) error {
	url := c.Query("url")
	if url == "" {
		return errors.New("missing url parameter")
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = util.ProcessImage(resp.Body, Conf.ReferenceDir)
	if err != nil {
		return err
	}
	return nil
}
