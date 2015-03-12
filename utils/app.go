package utils

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	CLIENT_PUBLIC_DIR = "/client/public"
	UPLOAD_DIR        = "/upload"

	DEFAULT_OUTPUT_DIR = "_sites"

	DEFAULT_LANG  = "en"
	DEFAULT_THEME = "willy" // @todo FIXME
)

func AppPublicDir() string {
	wd := viper.GetString("working_dir")
	if wd == "" {
		wd = WorkingDir()
	}
	return path.Join(wd, CLIENT_PUBLIC_DIR)
}

func AppUploadDir() string {
	return path.Join(AppPublicDir(), UPLOAD_DIR)
}

func AppUploadSiteDir(siteId string) string {
	return path.Join(AppUploadDir(), siteId)
}

func AppUploadSiteFilePath(siteId string, fileName string) string {
	return path.Join(AppUploadSiteDir(siteId), fileName)
}

func AppUploadSiteUrlPath(siteId string, fileName string) string {
	return path.Join(UPLOAD_DIR, siteId, fileName)
}

func AppEnsureUploadDir() {
	EnsureDirectory(AppUploadDir())
}

func AppEnsureSiteUploadDir(siteId string) {
	EnsureDirectory(AppUploadSiteDir(siteId))
}

func WorkingDir() string {
	result, _ := os.Getwd()

	return result
}
