package utils

import "errors"

/*
*

	Revert
	@description uses a condition function which will throw an error if the result is false, and nil if its true
	@example
	Revert(checkBool(), "it should be true!!!")
	func checkBool() bool {
		return false;
	}

*
*/

func Revert(result bool, msg string) error {

	if !result {
		return errors.New(msg)
	}
	return nil
}

func RevertNoParams(msg string) error {
	return errors.New(msg)
}
