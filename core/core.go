package core

import (
	"fmt"
	"os"
	"path"

	"github.com/aymerick/kowa/helpers"
	"github.com/spf13/viper"
)

const (
	UPLOAD_URL_PATH = "/upload"

	DEFAULT_LANG  = "en"
	DEFAULT_THEME = "willy" // @todo FIXME

	DEFAULT_BASEURL = "http://127.0.0.1"
)

func UploadDir() string {
	dir := viper.GetString("upload_dir")
	if dir == "" {
		panic("The upload_dir setting is mandatory")
	}

	return dir
}

func UploadSiteDir(siteId string) string {
	return path.Join(UploadDir(), siteId)
}

func UploadSiteFilePath(siteId string, fileName string) string {
	return path.Join(UploadSiteDir(siteId), fileName)
}

func UploadSiteUrlPath(siteId string, fileName string) string {
	return path.Join(UPLOAD_URL_PATH, siteId, fileName)
}

func EnsureUploadDir() {
	dir := UploadDir()

	parentDir, _ := path.Split(dir)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory %s does not exist. Your upload_dir setting may be incorrect.", parentDir))
	}

	helpers.EnsureDirectory(dir)
}

func EnsureSiteUploadDir(siteId string) {
	EnsureUploadDir()

	helpers.EnsureDirectory(UploadSiteDir(siteId))
}
