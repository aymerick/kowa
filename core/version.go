package core

import "fmt"

var (
	Version   string
	GitRev    string
	BuildDate string
)

func FormatVersion() string {
	var result string

	if Version != "" {
		result = fmt.Sprintf("kowa v%s [#%s] (%s)\n", Version, GitRev, BuildDate)
	} else {
		result = fmt.Sprintf("kowa (version unknown)\n")
	}

	return result
}
