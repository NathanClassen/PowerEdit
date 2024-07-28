package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func replaceQuotes(text string) string {
	// Replace straight double quotes with curly double quotes
	text = strings.ReplaceAll(text, "“", "\"")
	text = strings.ReplaceAll(text, "”", "\"")

	// Replace straight single quotes with curly single quotes
	text = strings.ReplaceAll(text, "‘", "'")
	text = strings.ReplaceAll(text, "’", "'")

	return text
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	// Read the file content
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Replace straight quotes with curly quotes
	modifiedContent := replaceQuotes(string(content))

	// Write the modified content to a new file
	newFilename := "modified_" + filename
	err = ioutil.WriteFile(newFilename, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Modified content written to %s\n", newFilename)
}
