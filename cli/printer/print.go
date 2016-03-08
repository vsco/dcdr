package printer

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).PrintfFunc()
	red   = color.New(color.FgRed).PrintfFunc()
)

func Say(format string, v ...interface{}) {
	green("[dcdr] ")
	fmt.Printf(format+"\n", v...)
}

func SayErr(format string, v ...interface{}) {
	red("[dcdr error] ")
	fmt.Printf(format+"\n", v...)
}
