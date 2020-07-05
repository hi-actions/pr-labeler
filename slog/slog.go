package slog

import (
	"os"

	"github.com/gookit/color"
)

// Infof log message
func Infof(fmt string, args ...interface{})  {
	color.Info.Printf(fmt + "\n", args...)
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

// Fatalf log message
func Fatalf(fmt string, args ...interface{})  {
	color.Error.Printf(fmt + "\n", args...)
	os.Exit(2)
}
