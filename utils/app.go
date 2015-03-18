package utils

import (
	"path"

	"github.com/spf13/viper"
)

const (
	APP_PUBLIC_DIR = "/public"
	UPLOAD_DIR     = "/upload"

	DEFAULT_OUTPUT_DIR = "_sites"

	DEFAULT_LANG  = "en"
	DEFAULT_THEME = "willy" // @todo FIXME

	DEFAULT_BASEURL = "http://127.0.0.1"
)

func AppDir() string {
	dir := viper.GetString("app_dir")
	if dir == "" {
		panic("The app_dir setting is mandatory")
	}

	return dir
}

func AppPublicDir() string {
	return path.Join(AppDir(), APP_PUBLIC_DIR)
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
