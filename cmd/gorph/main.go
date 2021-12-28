package main

import (
	"github.com/toudi/gorpheus/v1"
)

func main() {
	gorpheus.ScanDirectory("migrations")
}
