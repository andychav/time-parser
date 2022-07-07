package main

import (
	"os"
)

func main() {
	t := os.Args[1]

	p := new(parser)
	p.Buffer = t
	p.Init()
	p.Parse()
	p.Execute()
}
