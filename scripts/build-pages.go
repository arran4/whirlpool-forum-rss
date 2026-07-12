package main

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run build-pages.go <input-md> <output-html>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	md, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	// Add some basic HTML wrapper
	buf.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>Whirlpool Forum RSS</title>\n</head>\n<body>\n")

	if err := goldmark.Convert(md, &buf); err != nil {
		fmt.Printf("Error converting markdown: %v\n", err)
		os.Exit(1)
	}

	buf.WriteString("</body>\n</html>\n")

	err = os.MkdirAll(filepath.Dir(outputFile), 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}
