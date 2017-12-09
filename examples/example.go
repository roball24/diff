package main

import (
	"log"

	"github.com/roball24/go-diff"
)

func main() {
	var dc diff.Checker
	// not realistic, proof of concept
	if dc != nil {
		if err := dc.Run(); err != nil {
			log.Fatalf("DiffChecker Error: %s\n\n", err)
		}
	}
}
