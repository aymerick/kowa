package builder

import (
	"github.com/aymerick/raymond"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/aymerick/kowa/models"
)

func generateHTML(inputFormat string, input string) raymond.SafeString {
	var result raymond.SafeString

	if inputFormat == models.FormatMarkdown {
		html := blackfriday.MarkdownCommon([]byte(input))

		result = raymond.SafeString(bluemonday.UGCPolicy().SanitizeBytes(html))
	} else {
		sanitizePolicy := bluemonday.UGCPolicy()
		sanitizePolicy.AllowAttrs("style").OnElements("p", "span", "div") // I know this is bad
		sanitizePolicy.AllowAttrs("target").OnElements("a")

		result = raymond.SafeString(sanitizePolicy.Sanitize(input))
	}

	return result
}
