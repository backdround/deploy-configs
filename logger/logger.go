// logger describes type logger that log messages to system terminal
package logger

import (
	"fmt"
	"github.com/fatih/color"
)

type Logger interface {
	Title(title string)
	Success(message string)
	Warn(message string)
	Fail(message string)
	Log(message string)
}

type logger struct{}

func New() Logger {
	return &logger{}
}

func (l logger) Title(title string) {
	c := color.New(color.FgHiBlue, color.Bold)
	c.Printf("\n--- %v ---\n", title)
}

func (l logger) Success(message string) {
	c := color.New(color.FgHiGreen)
	c.Println(message)
}

func (l logger) Warn(message string) {
	c := color.New(color.FgYellow)
	c.Println(message)
}

func (l logger) Fail(message string) {
	c := color.New(color.FgHiRed, color.Bold)
	c.Println(message)
}

func (l logger) Log(message string) {
	fmt.Println(message)
}
