//Package errors contains custom error types for the Finance Manager application
package errors

import "fmt"

//type FmError is a custom error type that contains a REST response code and an error message
type FmError struct {
	Code int
	Msg  string
}

//Function Error overrides the default Error method. It returns the HTTP response code and then the message in a formatted manner
func (f *FmError) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.Msg)
}
