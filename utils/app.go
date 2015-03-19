package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	UPLOAD_URL_PATH = "/upload"

	DEFAULT_LANG  = "en"
	DEFAULT_THEME = "willy" // @todo FIXME

	DEFAULT_BASEURL = "http://127.0.0.1"
)

func AppUploadDir() string {
	dir := viper.GetString("upload_dir")
	if dir == "" {
		panic("The upload_dir setting is mandatory")
	}

	return dir
}

func AppUploadSiteDir(siteId string) string {
	return path.Join(AppUploadDir(), siteId)
}

func AppUploadSiteFilePath(siteId string, fileName string) string {
	return path.Join(AppUploadSiteDir(siteId), fileName)
}

func AppUploadSiteUrlPath(siteId string, fileName string) string {
	return path.Join(UPLOAD_URL_PATH, siteId, fileName)
}

func AppEnsureUploadDir() {
	dir := AppUploadDir()

	parentDir, _ := path.Split(dir)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory %s does not exist. Your upload_dir setting may be incorrect.", parentDir))
	}

	EnsureDirectory(dir)
}

func AppEnsureSiteUploadDir(siteId string) {
	AppEnsureUploadDir()

	EnsureDirectory(AppUploadSiteDir(siteId))
}
