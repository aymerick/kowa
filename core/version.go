package core

import "fmt"

var (
	// Version is kowa version
	Version string

	// GitRev is kowa GIT revision
	GitRev string

	// BuildDate is kowa build date
	BuildDate string
)

// FormatVersion returns the human readable kowa version
func FormatVersion() string {
	var result string

	if Version != "" {
		result = fmt.Sprintf("kowa v%s [#%s] (%s)\n", Version, GitRev, BuildDate)
	} else {
		result = fmt.Sprintf("kowa (version unknown)\n")
	}

	return result
}
