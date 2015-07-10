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
	DEFAULT_TZ    = "Europe/Paris"
	DEFAULT_THEME = "willy"

	DEFAULT_BASEURL = "http://127.0.0.1"
)

// DefaultDomain returns default domain, or an empty string if no domain found in settings.
func DefaultDomain() string {
	domains := viper.GetStringSlice("service_domains")
	if len(domains) == 0 {
		return ""
	}

	return domains[0]
}

// ValidDomain returns true if domain is valid.
func ValidDomain(domain string) bool {
	for _, serviceDomain := range viper.GetStringSlice("service_domains") {
		if serviceDomain == domain {
			return true
		}
	}

	return false
}

// BaseUrl computes a base url for given site id.
func BaseUrl(siteId string) string {
	if domain := DefaultDomain(); domain != "" {
		return BaseUrlForDomain(siteId, domain)
	}

	return fmt.Sprintf("%s:%d/%s", DEFAULT_BASEURL, viper.GetInt("serve_output_port"), siteId)
}

// BaseUrlForDomain computes a base url for given site id and domain.
func BaseUrlForDomain(siteId string, domain string) string {
	return fmt.Sprintf("http://%s.%s", siteId, domain)
}

// BaseUrlForCustomDomain computes a base url for given custom domain.
func BaseUrlForCustomDomain(domain string) string {
	return fmt.Sprintf("http://%s", domain)
}

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
