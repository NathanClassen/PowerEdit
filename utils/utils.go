package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)



// readWords reads the words from a file and returns them as a slice.
func ReadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}

// printSurroundingWords prints the surrounding words from each file around the discrepancy.
func PrintSurroundingWords(words1, words2 []string, toEditIndex, sourceIndex int, file1, file2 string) {
	istart := max(0, toEditIndex-10)
	iend := min(len(words1), toEditIndex+11)
	jstart := max(0, sourceIndex-10)
	jend := min(len(words1), sourceIndex+11)

	fmt.Printf("Discrepancy found:\n")
	fmt.Printf("%25s: %s *%s* %s\n", file1, strings.Join(words1[istart:toEditIndex], " "), words1[toEditIndex], words1[toEditIndex+1:iend])
	fmt.Printf("%25s: %s *%s* %s\n", file2, strings.Join(words2[jstart:sourceIndex], " "), words2[sourceIndex], words2[sourceIndex+1:jend])
}

// max returns the maximum of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// updateFile updates the file content with the specified words.
func UpdateFile(filename string, words []string) error {
	content := strings.Join(words, " ")
	return os.WriteFile(filename, []byte(content), 0644)
}

func ParseDigits(input string) ([]int, error) {
	var digits []int
	for _, char := range input {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, err
		}
		digits = append(digits, digit)
	}
	return digits, nil
}

func ReplaceQuotes(text string) string {
	// Replace straight double quotes with curly double quotes
	text = strings.ReplaceAll(text, "“", "\"")
	text = strings.ReplaceAll(text, "”", "\"")

	// Replace straight single quotes with curly single quotes
	text = strings.ReplaceAll(text, "‘", "'")
	text = strings.ReplaceAll(text, "’", "'")

	return text
}

func CleanWord(s string) string {
	// Create a new slice to hold the filtered runes
	var result []rune
	// Iterate over each rune in the input string
	for _, r := range s {
		// Check if the rune is a letter
		if unicode.IsLetter(r) {
			// If it is a letter, append it to the result slice
			result = append(result, r)
		}
	}
	// Convert the result slice back to a string and return it
	return string(result)
}
