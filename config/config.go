package config

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	MultipartDir string
	Base64Dir    string
	ReferenceDir string
}

var Conf config

func init() {
	viper.AddConfigPath("./config")
	if os.Getenv("ENVIRONMENT") == "DEV" {
		viper.SetConfigName("config-local")
	} else {
		viper.SetConfigName("config")
	}
	viper.ReadInConfig()
	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("cannot unmarshal config into struct, %v", err)
	}
}

func InitTestConfig() string {
	dir, _ := ioutil.TempDir("", "thumbnail_server")
	multipartDir, _ := ioutil.TempDir(dir, "multipart")
	base64Dir, _ := ioutil.TempDir(dir, "base64")
	referenceDir, _ := ioutil.TempDir(dir, "reference")

	Conf = config{
		multipartDir,
		base64Dir,
		referenceDir,
	}
	return dir
}
