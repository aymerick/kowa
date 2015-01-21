package builder

import "log"

type ErrorCollector struct {
	Errors map[string][]error
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		Errors: make(map[string][]error),
	}
}

// add a new error for given step
func (collector *ErrorCollector) AddError(step string, err error) {
	collector.Errors[step] = append(collector.Errors[step], err)
}

// dump all errors
func (collector *ErrorCollector) Dump() {
	for step, errors := range collector.Errors {
		log.Printf("=== step: %s", step)

		for _, err := range errors {
			log.Printf(err.Error())
		}
	}
}
