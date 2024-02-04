//Package fmlogger contains custom logging methods
package fmlogger

import "fmt"

const entry_msg = "[ENTER %s]\n"
const exit_msg = "[EXIT %s]\n"
const err_msg = "[%s] %s:\n%v"
const info_msg = "[%s] %s\n"
const info_obj_msg = "[%s] %s:\n%v"

//Function Enter returns a formated string used to declare where a method begins execution
func Enter(method string) {
	fmt.Printf(entry_msg, method)
}

//Function Error returns a formated string used to log a given error along with a custom error message and declaring which method the error occured in
func Error(method string, msg string, err error) {
	fmt.Printf(err_msg+"\n", method, msg, err)
}

//Function Exit returns a formated string used to declare where a method ends execution
func Exit(method string) {
	fmt.Printf(exit_msg, method)
}

//Function ExitError returns a formated string used to combine the Exit and Error functions together
func ExitError(method string, msg string, err error) {
	Error(method, msg, err)
	Exit(method)
}

//Fucntion Info returns a formatted string containing a custom message and the method that the message is coming from
func Info(method string, msg string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(info_msg, method, msg), args...)
}

//Function InfoObj returns a formatted string containing a custom message, the method that the message is coming from, and one object of any type
func InfoObj(method string, msg string, obj any) {
	fmt.Printf(info_obj_msg, method, msg, obj)
}
