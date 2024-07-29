package utils

import (
	"fmt"
	"strings"
)

type Errors struct {
	errs []error
}

func (e *Errors) Append(err error) {
	e.errs = append(e.errs, err)
}

func (e *Errors) Error() error {
	if len(e.errs) == 0 {
		return nil
	}

	strs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		strs = append(strs, err.Error())
	}

	return fmt.Errorf(strings.Join(strs, "; "))
}
