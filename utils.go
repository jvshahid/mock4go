package api

import (
	"fmt"
)

var verbose = false

func SetVerbosity(isVerbose bool) {
	verbose = isVerbose
}

func Log(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg, args...)
	}
}
