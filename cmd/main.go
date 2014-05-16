// Copyright 2010 Andrey Mirtchovski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go distribution's
// LICENSE file.
package main

// Test program for the iterative walker
import (
	"fmt"
	"log"
	"os"

	"github.com/mirtchovski/walk"
)

func main() {
	for _, p := range os.Args[1:] {
		err := walk.Walk(p, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			} else {
				fmt.Println(path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
