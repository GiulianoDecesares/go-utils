package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadStdInAsInt() (int, error) {
	var err error
	var userOption string
	var result int

	if userOption, err = ReadStdInAsString(); err == nil {
		result, err = strconv.Atoi(userOption)
	}

	return result, err
}

func ReadStdInAsString() (string, error) {
	var err error = nil
	var userOption string = ""

	if userOption, err = bufio.NewReader(os.Stdin).ReadString('\n'); err == nil {
		userOption = strings.TrimSuffix(userOption, "\n")
		userOption = strings.TrimSuffix(userOption, "\r")
	}

	return userOption, err
}
