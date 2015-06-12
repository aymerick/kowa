package builder

import (
	"reflect"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

// Build FuncMap for template
func (builder *SiteBuilder) helpers() map[string]interface{} {
	return map[string]interface{}{
		"urlFor":     builder.UrlFor,
		"startsWith": builder.StartsWith,
		"t":          builder.Translate,
		"mod":        builder.Mod,
		"modBool":    builder.ModBool,
	}
}

// UrlFor returns an URL to an internal page.
func (builder *SiteBuilder) UrlFor(dest string) string {
	var result string

	switch dest {
	case KIND_ACTIVITIES, KIND_MEMBERS, KIND_CONTACT, KIND_HOMEPAGE:
		// find uniq node
		nodes := builder.nodeBuilder(dest).Nodes()
		if len(nodes) == 0 {
			panic("No node loaded yet")
		}

		if len(nodes) > 1 {
			panic("That method logic is broken, fix it!")
		}

		result = nodes[0].Url

	case KIND_POSTS, KIND_EVENTS:
		// find correct node
		nodes := builder.nodeBuilder(dest).Nodes()

		for _, node := range nodes {
			if node.Kind == dest {
				result = node.Url
				break
			}
		}

	default:
		// KIND_PAGE, KIND_POST
		panic("Internal link kind not supported yet")
	}

	return result
}

// StartsWith returns true if check has prefix.
func (builder *SiteBuilder) StartsWith(check string, prefix string) bool {
	return strings.HasPrefix(check, prefix)
}

// Translate translates given sentence.
func (builder *SiteBuilder) Translate(sentence string) string {
	T := i18n.MustTfunc(builder.site.Lang)

	return T(sentence)
}

// @todo Replace interface{}
func (builder *SiteBuilder) Mod(a, b interface{}) int64 {
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

func (builder *SiteBuilder) ModBool(a, b interface{}) bool {
	return builder.Mod(a, b) == int64(0)
}
