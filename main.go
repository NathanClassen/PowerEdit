package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// readWords reads the words from a file and returns them as a slice.
func readWords(filename string) ([]string, error) {
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
func printSurroundingWords(words1, words2 []string, index int) {
    start := max(0, index-5)
    end := min(len(words1), index+6)

    fmt.Printf("Discrepancy found:\n")
    fmt.Printf("File 1: %s\n", strings.Join(words1[start:end], " "))
    fmt.Printf("File 2: %s\n", strings.Join(words2[start:end], " "))
}

// max returns the maximum of two integers.
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// min returns the minimum of two integers.
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: compare <file1> <file2>")
        return
    }

    file1 := os.Args[1]
    file2 := os.Args[2]

    words1, err := readWords(file1)
    if err != nil {
        fmt.Println("Error reading file1:", err)
        return
    }

    words2, err := readWords(file2)
    if err != nil {
        fmt.Println("Error reading file2:", err)
        return
    }

    minLength := min(len(words1), len(words2))
    for i := 0; i < minLength; i++ {
        if words1[i] != words2[i] {
            printSurroundingWords(words1, words2, i)
            return
        }
    }

    if len(words1) != len(words2) {
        fmt.Println("Files have different lengths.")
    } else {
        fmt.Println("Files are identical.")
    }
}
