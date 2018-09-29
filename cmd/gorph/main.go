package main

import (
	"github.com/toudi/gorpheus"
)

func main() {
	gorpheus.ScanDirectory("migrations")
}
