package main

import (
	"fmt"
	"log"

	"github.com/roball24/diff"
)

func main() {
	equal, err := diff.DefaultFileCheck("file.txt", "equal.txt", true) // file1, file2, verbose
	if err != nil {
		log.Fatalf("FileCheck Error: %s\n\n", err)
	}
	fmt.Printf("diffcheck...%t\n", equal)
}
