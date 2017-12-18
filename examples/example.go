package main

import (
	"bytes"
	"log"

	"github.com/roball24/go-diff"
)

func main() {
	var dc diff.Checker
	// not realistic, proof of concept
	if dc != nil {
		dc.AddLineCompare(bytes.Equal, diff.ErrorFail)
		if err := dc.Run(); err != nil {
			log.Fatalf("DiffChecker Error: %s\n\n", err)
		}
	}
}
