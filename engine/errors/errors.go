package errors

import "fmt"

type Op string

func Error(op Op, desc string, err error) error {
	if err == nil {
		return fmt.Errorf("%s: %s", op, desc)
	}
	return fmt.Errorf("%s: %s: %w", op, desc, err)
}
