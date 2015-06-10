package raymond

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"
)

// Template represents a handlebars template.
type Template struct {
	source   string
	program  *ast.Program
	helpers  map[string]reflect.Value
	partials map[string]*partial
}

// newTemplate instanciate a new template without parsing it
func newTemplate(source string) *Template {
	return &Template{
		source:   source,
		helpers:  make(map[string]reflect.Value),
		partials: make(map[string]*partial),
	}
}

// Parse instanciates a template by parsing given source.
func Parse(source string) (*Template, error) {
	tpl := newTemplate(source)

	// parse template
	if err := tpl.parse(); err != nil {
		return nil, err
	}

	return tpl, nil
}

// MustParse instanciates a template by parsing given source. It panics on error.
func MustParse(source string) *Template {
	result, err := Parse(source)
	if err != nil {
		panic(err)
	}
	return result
}

// parse parses the template
//
// It can be called several times, the parsing will be done only once.
func (tpl *Template) parse() error {
	if tpl.program == nil {
		var err error

		tpl.program, err = parser.Parse(tpl.source)
		if err != nil {
			return err
		}
	}

	return nil
}

// Clone returns a copy of that template.
func (tpl *Template) Clone() *Template {
	result := newTemplate(tpl.source)

	result.program = tpl.program

	for name, helper := range tpl.helpers {
		result.RegisterHelper(name, helper)
	}

	for name, partial := range tpl.partials {
		result.addPartial(name, partial.source, partial.tpl)
	}

	return result
}

// RegisterHelper registers a helper for that template.
func (tpl *Template) RegisterHelper(name string, helper interface{}) {
	if tpl.helpers[name] != zero {
		panic(fmt.Sprintf("Helper %s already registered", name))
	}

	val := reflect.ValueOf(helper)
	ensureValidHelper(name, val)

	tpl.helpers[name] = val
}

// RegisterHelpers registers several helpers for that template.
func (tpl *Template) RegisterHelpers(helpers map[string]interface{}) {
	for name, helper := range helpers {
		tpl.RegisterHelper(name, helper)
	}
}

func (tpl *Template) addPartial(name string, source string, template *Template) {
	tpl.partials[name] = newPartial(name, source, template)
}

// RegisterPartial registers a partial for that template.
func (tpl *Template) RegisterPartial(name string, source string) {
	if tpl.partials[name] != nil {
		panic(fmt.Sprintf("Partial %s already registered", name))
	}

	tpl.addPartial(name, source, nil)
}

// RegisterPartials registers several partials for that template.
func (tpl *Template) RegisterPartials(partials map[string]string) {
	for name, partial := range partials {
		tpl.RegisterPartial(name, partial)
	}
}

// RegisterPartial registers an already parsed partial for that template.
func (tpl *Template) RegisterPartialTemplate(name string, template *Template) {
	if tpl.partials[name] != nil {
		panic(fmt.Sprintf("Partial %s already registered", name))
	}

	tpl.addPartial(name, "", template)
}

// Exec evaluates template with given context.
func (tpl *Template) Exec(ctx interface{}) (result string, err error) {
	return tpl.ExecWith(ctx, nil)
}

// MustExec evaluates template with given context. It panics on error.
func (tpl *Template) MustExec(ctx interface{}) string {
	result, err := tpl.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

// ExecWith evaluates template with given context and private data frame.
func (tpl *Template) ExecWith(ctx interface{}, privData *DataFrame) (result string, err error) {
	defer errRecover(&err)

	// parses template if necessary
	err = tpl.parse()
	if err != nil {
		return
	}

	// setup visitor
	v := newEvalVisitor(tpl, ctx, privData)

	// visit AST
	result, _ = tpl.program.Accept(v).(string)

	// named return values
	return
}

// errRecover recovers evaluation panic
func errRecover(errp *error) {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case runtime.Error:
			panic(e)
		case error:
			*errp = err
		default:
			panic(e)
		}
	}
}

// PrintAST returns string representation of parsed template.
func (tpl *Template) PrintAST() string {
	if err := tpl.parse(); err != nil {
		return fmt.Sprintf("PARSER ERROR: %s", err)
	}

	return ast.Print(tpl.program)
}
