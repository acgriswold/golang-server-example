package handler

import (
	"fmt"
	"os"

	"errors"
)

func Check(err error, msg string, passError bool) {
	if err != nil {
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

func CheckFileError(err error, msg string) {
	if !errors.Is(err, os.ErrNotExist) {
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}
