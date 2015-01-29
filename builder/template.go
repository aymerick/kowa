package builder

import (
	"errors"
	"reflect"

	"html/template"
)

var templateFuncMap template.FuncMap

func init() {
	templateFuncMap = template.FuncMap{
		"mod":     Mod,
		"modBool": ModBool,
	}
}

// borrowed from https://github.com/spf13/hugo/blob/master/tpl/template.go
func Mod(a, b interface{}) (int64, error) {
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

// borrowed from https://github.com/spf13/hugo/blob/master/tpl/template.go
func ModBool(a, b interface{}) (bool, error) {
	res, err := Mod(a, b)
	if err != nil {
		return false, err
	}
	return res == int64(0), nil
}
