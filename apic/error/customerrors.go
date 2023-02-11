package errors

import "fmt"

type BadParamsError struct {
}

func (e *BadParamsError) Error() string {
	return fmt.Sprintf("Invalid params:" )
}

func NewBadParamsError() *BadParamsError {
	return &BadParamsError{}
}
