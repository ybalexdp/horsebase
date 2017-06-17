package main

import (
	"fmt"
	"io"

	"github.com/mitchellh/colorstring"
)

func PrintError(w io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]\n", format)
	fmt.Fprintf(w, colorstring.Color(fmt.Sprintf(format, args...)))
}
