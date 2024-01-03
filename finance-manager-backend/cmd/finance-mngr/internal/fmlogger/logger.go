package fmlogger

import "fmt"

const entry_msg = "[ENTER %s]\n"
const exit_msg = "[EXIT %s]\n"
const err_msg = "[%s] %s:\n%v"
const info_msg = "[%s] %s\n"
const info_obj_msg = "[%s] %s:\n%v"

func Enter(method string) {
	fmt.Printf(entry_msg, method)
}

func Error(method string, msg string, err error) {
	fmt.Printf(err_msg+"\n", method, msg, err)
}

func Exit(method string) {
	fmt.Printf(exit_msg, method)
}

func ExitError(method string, msg string, err error) {
	Error(method, msg, err)
	Exit(method)
}

func Info(method string, msg string) {
	fmt.Printf(info_msg, method, msg)
}

func InfoObj(method string, msg string, obj any) {
	fmt.Printf(info_obj_msg, method, msg, obj)
}
