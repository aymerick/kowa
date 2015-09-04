package core

import (
	"fmt"
	"os"
	"path"

	"github.com/aymerick/kowa/helpers"
	"github.com/spf13/viper"
)

const (
	// DefaultLang is default language
	DefaultLang = "en"

	// DefaultTZ is default timezone
	DefaultTZ = "Europe/Paris"

	// DefaultTheme is default theme
	DefaultTheme = "willy"

	uploadURLPath  = "/upload"
	defaultBaseURL = "http://127.0.0.1"
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
func BaseUrl(siteID string) string {
	if domain := DefaultDomain(); domain != "" {
		return BaseUrlForDomain(siteID, domain)
	}

	return fmt.Sprintf("%s:%d/%s", defaultBaseURL, viper.GetInt("serve_output_port"), siteID)
}

// BaseUrlForDomain computes a base url for given site id and domain.
func BaseUrlForDomain(siteID string, domain string) string {
	return fmt.Sprintf("http://%s.%s", siteID, domain)
}

// BaseUrlForCustomDomain computes a base url for given custom domain.
func BaseUrlForCustomDomain(domain string) string {
	return fmt.Sprintf("http://%s", domain)
}

// UploadDir returns the main upload directory path
func UploadDir() string {
	dir := viper.GetString("upload_dir")
	if dir == "" {
		panic("The upload_dir setting is mandatory")
	}

	return dir
}

// UploadSiteDir returns the upload directory path for given site
func UploadSiteDir(siteID string) string {
	return path.Join(UploadDir(), siteID)
}

// UploadSiteFilePath returns path to uploaded file for given site
func UploadSiteFilePath(siteID string, fileName string) string {
	return path.Join(UploadSiteDir(siteID), fileName)
}

// UploadSiteUrlPath returns relative URL to uploaded file for given site
func UploadSiteUrlPath(siteID string, fileName string) string {
	return path.Join(uploadURLPath, siteID, fileName)
}

// EnsureUploadDir ensures that main upload directory exists
func EnsureUploadDir() {
	dir := UploadDir()

	parentDir, _ := path.Split(dir)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory %s does not exist. Your upload_dir setting may be incorrect.", parentDir))
	}

	helpers.EnsureDirectory(dir)
}

// EnsureUploadDir ensures that the upload directory for given site exists
func EnsureSiteUploadDir(siteID string) {
	EnsureUploadDir()

	helpers.EnsureDirectory(UploadSiteDir(siteID))
}
