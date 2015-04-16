package builder

import (
	"errors"
	"html/template"
	"reflect"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/russross/blackfriday"

	"github.com/aymerick/kowa/models"
)

func generateHTML(inputFormat string, input string) template.HTML {
	var result template.HTML

	if inputFormat == models.FORMAT_MARKDOWN {
		html := blackfriday.MarkdownCommon([]byte(input))

		result = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))
	} else {
		sanitizePolicy := bluemonday.UGCPolicy()
		sanitizePolicy.AllowAttrs("style").OnElements("p", "span", "div") // I know this is bad
		sanitizePolicy.AllowAttrs("target").OnElements("a")

		result = template.HTML(sanitizePolicy.Sanitize(input))
	}

	return result
}

//
// Func Map
//

// Build FuncMap for template
func (builder *SiteBuilder) FuncMap() template.FuncMap {
	return template.FuncMap{
		"urlFor":     builder.UrlFor,
		"startsWith": builder.StartsWith,
		"t":          builder.Translate,
		"mod":        builder.Mod,
		"modBool":    builder.ModBool,
	}
}

// Returns an URL to an internal page
func (builder *SiteBuilder) UrlFor(dest interface{}) (string, error) {
	destValue := reflect.ValueOf(dest)

	var destStr string

	switch destValue.Kind() {
	case reflect.String:
		destStr = destValue.String()
	default:
		return "", errors.New("urlFor operator needs a string argument")
	}

	var result string
	var err error

	switch destStr {
	case KIND_ACTIVITIES, KIND_MEMBERS, KIND_CONTACT, KIND_HOMEPAGE:
		// find uniq node
		nodes := builder.nodeBuilder(destStr).Nodes()
		if len(nodes) == 0 {
			err = errors.New("No node loaded yet")
		} else if len(nodes) > 1 {
			err = errors.New("That method logic is broken, fix it!")
		} else {
			result = nodes[0].Url
		}

	case KIND_POSTS, KIND_EVENTS:
		// find correct node
		nodes := builder.nodeBuilder(destStr).Nodes()

		for _, node := range nodes {
			if node.Kind == destStr {
				result = node.Url
				break
			}
		}

	default:
		// KIND_PAGE, KIND_POST
		err = errors.New("Internal link kind not supported yet")
	}

	return result, err
}

// Returns an URL to an internal page
func (builder *SiteBuilder) StartsWith(check interface{}, prefix interface{}) (bool, error) {
	checkValue := reflect.ValueOf(check)
	prefixValue := reflect.ValueOf(prefix)

	var checkStr string
	var prefixStr string

	switch checkValue.Kind() {
	case reflect.String:
		checkStr = checkValue.String()
	default:
		return false, errors.New("startsWith operator needs string arguments")
	}

	switch prefixValue.Kind() {
	case reflect.String:
		prefixStr = prefixValue.String()
	default:
		return false, errors.New("startsWith operator needs string arguments")
	}

	return strings.HasPrefix(checkStr, prefixStr), nil
}

func (builder *SiteBuilder) Translate(sentence interface{}) (string, error) {
	sentenceValue := reflect.ValueOf(sentence)

	if sentenceValue.Kind() != reflect.String {
		return "", errors.New("Translate operator needs a string argument")
	}

	sentenceStr := sentenceValue.String()

	T := i18n.MustTfunc(builder.site.Lang)

	return T(sentenceStr), nil
}

// Borrowed from https://github.com/spf13/hugo/blob/master/tpl/template.go
func (builder *SiteBuilder) Mod(a, b interface{}) (int64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bi = bv.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("The number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// Borrowed from https://github.com/spf13/hugo/blob/master/tpl/template.go
func (builder *SiteBuilder) ModBool(a, b interface{}) (bool, error) {
	res, err := builder.Mod(a, b)
	if err != nil {
		return false, err
	}
	return res == int64(0), nil
}
