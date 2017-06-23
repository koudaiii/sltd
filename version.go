package main

import (
	"fmt"
)

var (
	Version  string
	Revision string
)

func printVersion() {
	fmt.Println("sltd version " + Version + ", build " + Revision)
}
