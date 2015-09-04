package builder

import (
	"reflect"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

// Build FuncMap for template
func (site *SiteBuilder) helpers() map[string]interface{} {
	return map[string]interface{}{
		"urlFor": site.UrlFor,
		"t":      site.Translate,

		"startsWith": StartsWith,
		"mod":        Mod,
		"modBool":    ModBool,
	}
}

// UrlFor returns an URL to an internal page.
func (site *SiteBuilder) UrlFor(dest string) string {
	var result string

	switch dest {
	case kindActivities, kindMembers, kindContact, kindHomepage:
		// find uniq node
		nodes := site.nodeBuilder(dest).Nodes()
		if len(nodes) == 1 {
			result = nodes[0].Url
		}

	case kindPosts, kindEvents:
		// find correct node
		nodes := site.nodeBuilder(dest).Nodes()

		for _, node := range nodes {
			if node.Kind == dest {
				result = node.Url
				break
			}
		}

	default:
		// kindPage, kindPost
		panic("Internal link kind not supported yet")
	}

	return result
}

// Translate translates given sentence.
func (site *SiteBuilder) Translate(sentence string) string {
	T := i18n.MustTfunc(site.site.Lang)

	return T(sentence)
}

// StartsWith returns true if check has prefix.
func StartsWith(check string, prefix string) bool {
	return strings.HasPrefix(check, prefix)
}

// Mod returns modulo of given int values
// @todo Replace interface{}
func Mod(a, b interface{}) int64 {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
	default:
		panic("Modulo operator can't be used with non integer value")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bi = bv.Int()
	default:
		panic("Modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		panic("The number can't be divided by zero at modulo operation")
	}

	return ai % bi
}

// ModBool returns true if the modulo of given integers is zero
func ModBool(a, b interface{}) bool {
	return Mod(a, b) == int64(0)
}
