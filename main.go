package main

import "os"

func main() {
	hb := &Horsebase{}
	hb = hb.New()
	os.Exit(hb.Run(os.Args))
}
