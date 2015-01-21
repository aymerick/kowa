package builder

import (
	"fmt"
	"log"
)

type ErrorCollector struct {
	Errors   map[string][]error
	ErrorsNb int
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		Errors:   make(map[string][]error),
		ErrorsNb: 0,
	}
}

// add a new error for given step
func (collector *ErrorCollector) AddError(step string, err error) {
	collector.Errors[step] = append(collector.Errors[step], err)

	collector.ErrorsNb += 1
}

// dump all errors
func (collector *ErrorCollector) Dump() {
	if collector.ErrorsNb > 0 {
		log.Printf("[ERR] Built with %d error(s)", collector.ErrorsNb)

		errNb := 1

		for step, errors := range collector.Errors {
			log.Printf("[ERR] %s:", step)

			for _, err := range errors {
				log.Printf(fmt.Sprintf("[ERR]   %d. %v", errNb, err.Error()))
				errNb += 1
			}
		}
	}
}
