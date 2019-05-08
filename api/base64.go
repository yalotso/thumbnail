package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-ozzo/ozzo-routing"
	. "github.com/yalotso/thumbnail/config"
	"github.com/yalotso/thumbnail/util"
	"strings"
)

func Base64(c *routing.Context) error {
	if c.Request.Body == nil {
		return errors.New("missing request body")
	}
	dec := json.NewDecoder(c.Request.Body)
	var base64string string
	err := dec.Decode(&base64string)
	if err != nil {
		return err
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64string))
	err = util.ProcessImage(reader, Conf.Base64Dir)
	if err != nil {
		return err
	}
	return nil
}
