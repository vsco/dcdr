package printer

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).PrintfFunc()
	red   = color.New(color.FgRed).PrintfFunc()
)

func Info(format string, v ...interface{}) string {
	green("[dcdr] ")
	return fmt.Sprintf(fmt.Sprintf("%s\n", format), v...)
}

func Err(format string, v ...interface{}) string {
	red("[dcdr error] ")
	return fmt.Sprintf(fmt.Sprintf("%s\n", format), v...)
}

func Say(format string, v ...interface{}) {
	fmt.Printf(Info(format, v...))
}

func SayErr(format string, v ...interface{}) {
	fmt.Printf(Err(format, v...))
}

func Log(format string, v ...interface{}) {
	log.Printf(Info(format, v...))
}

func LogErr(format string, v ...interface{}) {
	log.Printf(Err(format, v...))
}
