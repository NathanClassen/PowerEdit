package main

import (
	"fmt"
	"os"
)

func replaceQuotes(text string) string {
	var result []rune
	inDoubleQuote := false
	inSingleQuote := false

	for _, char := range text {
		switch char {
		case '"':
			if inDoubleQuote {
				result = append(result, '”')
			} else {
				result = append(result, '“')
			}
			inDoubleQuote = !inDoubleQuote
		case '\'':
			if inSingleQuote {
				result = append(result, '’')
			} else {
				result = append(result, '‘')
			}
			inSingleQuote = !inSingleQuote
		default:
			// if unicode.IsSpace(char) || unicode.IsPunct(char) {
			// 	if inDoubleQuote {
			// 		inDoubleQuote = false
			// 		result = append(result, '”')
			// 	}
			// 	if inSingleQuote {
			// 		inSingleQuote = false
			// 		result = append(result, '’')
			// 	}
			// }
			result = append(result, char)
		}
	}

	return string(result)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Replace straight quotes with curly quotes
	modifiedContent := replaceQuotes(string(content))

	// Write the modified content to a new file
	newFilename := "modified_" + filename
	err = os.WriteFile(newFilename, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Modified content written to %s\n", newFilename)
}
