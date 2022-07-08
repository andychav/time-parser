package main

import (
	"fmt"
	"os"
)

func main() {
	t := os.Args[1]

	// Return minutes using parser from grammar
	p := new(parser)
	p.Buffer = t
	p.Init()
	p.Parse()
	p.Execute()

	// Return minutes using Regex
	fmt.Println("Using regex: ", getMinAfterMidnight(t))
}
