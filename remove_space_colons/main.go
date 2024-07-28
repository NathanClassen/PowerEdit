package main

import (
	"fmt"
	"os"
	"regexp"
)

// removeSpaces removes spaces between words and colons/semicolons.
func removeSpaces(input string) string {
	re := regexp.MustCompile(`\s+([;:'"?!.,])`)
	return re.ReplaceAllString(input, "$1")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: remove_spaces <filename>")
		return
	}

	filename := os.Args[1]

	// Read the entire file into memory.
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Remove spaces between words and colons/semicolons.
	output := removeSpaces(string(input))

	// Write the modified content back to the file.
	err = os.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File updated successfully.")
}
