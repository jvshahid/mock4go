package main

import (
	"fmt"
	"gomock"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s packages\n", os.Args[0])
		os.Exit(1)
	}

	for _, pkg := range os.Args[1:] {
		files, _ := gomock.GetFiles(pkg)
		fmt.Printf("package %s contains: %s\n", pkg, strings.Join(files, ","))
		for _, file := range files {
			content, err := gomock.InstrumentFile(pkg + "/" + file)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Content of %s:\n%s\n", file, content)
		}
	}
}
