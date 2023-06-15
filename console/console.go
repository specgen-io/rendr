package console

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	NoneLevel    LogLevel = 0
	VerboseLevel LogLevel = 1
)

var Level = NoneLevel

func Verbose(format string, args ...interface{}) {
	if Level >= VerboseLevel {
		fmt.Fprintf(os.Stdout, format, args...)
		fmt.Fprintln(os.Stdout)
	}
}

func Error(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
