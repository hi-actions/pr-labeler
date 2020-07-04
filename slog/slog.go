package slog

import (
	"os"

	"github.com/gookit/color"
)

// Infof log message
func Infof(fmt string, args ...interface{})  {
	color.Info.Printf(fmt, args...)
}

// Error log message
func Error(args ...interface{})  {
	color.Error.Println(args...)
}

// Fatal log message
func Fatal(args ...interface{})  {
	color.Error.Println(args...)
	os.Exit(2)
}
