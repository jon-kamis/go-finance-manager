package errors

import "fmt"

type FmError struct {
	Code int
	Msg  string
}

func (f *FmError) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.Msg)
}
