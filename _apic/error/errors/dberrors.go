package errors

import "fmt"

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "no entries found"
}

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{}
}


type ScanningIdError struct {
	Id int
}

func (e *ScanningIdError) Error() string {
	return fmt.Sprintf("error scanning bank with ID %v", e.Id)
}

func NewScanningIdError(id int) *ScanningIdError {
	return &ScanningIdError{Id: id}
}


type IdNotFoundError struct {
	Id int
}

func (e *IdNotFoundError) Error() string {
	return fmt.Sprintf("id: %v not found", e.Id)
}

func NewIdNotFoundError(id int) *IdNotFoundError {
	return &IdNotFoundError{Id: id}
}


type CreatingBankError struct {
}

func (e *CreatingBankError) Error() string {
	return fmt.Sprintf("error creating bank")
}

func NewCreatingBankError() *CreatingBankError {
	return &CreatingBankError{}
}

type DeletingBankError struct {
}

func (e *DeletingBankError) Error() string {
	return "error creating bank"
}

func NewDeletingBankError() *DeletingBankError {
	return &DeletingBankError{}
}
