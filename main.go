package main

import "os"

var version string

func main() {
	hb := &Horsebase{}
	hb = hb.New()
	os.Exit(hb.Run(os.Args))
}
