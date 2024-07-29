
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
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
    func printSurroundingWords(words1, words2 []string, toEditIndex, sourceIndex int, file1, file2 string) {
        istart := max(0, toEditIndex-10)
        iend := min(len(words1), toEditIndex+11)
        jstart := max(0, sourceIndex-10)
        jend := min(len(words1), sourceIndex+11)

        fmt.Printf("Discrepancy found:\n")
        fmt.Printf("%25s: %s *%s* %s\n", file1, strings.Join(words1[istart:toEditIndex], " "), words1[toEditIndex], words1[toEditIndex+1:iend])
        fmt.Printf("%25s: %s *%s* %s\n", file2, strings.Join(words2[jstart:sourceIndex], " "), words2[sourceIndex], words2[sourceIndex+1:jend])
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

    // updateFile updates the file content with the specified words.
    func updateFile(filename string, words []string) error {
        content := strings.Join(words, " ")
        return os.WriteFile(filename, []byte(content), 0644)
    }

    func main() {

        /* ************************************************************************
            GET OPTIONS SET
        ************************************************************************ */
        if len(os.Args) < 3 {
            fmt.Println("Usage: compare <file to edit> <'source file' upon which edits are based>")
            return
        }

        idx := 0
        jdx := 0

        file1 := os.Args[1]     // file to edit
        file2 := os.Args[2]     // source file

        idx, err := strconv.Atoi(os.Args[3])
        if err != nil {
            fmt.Printf("Couldn't parse idx, using 0")
            idx = 0
        }

        jdx, err = strconv.Atoi(os.Args[4])
        if err != nil {
            fmt.Printf("Couldn't parse jdx, using 0")
            jdx = 0
        }

        /* ************************************************************************
            READ IN WORDS OF FILE 1
        ************************************************************************ */
      
        words1, err := readWords(file1)
        if err != nil {
            fmt.Println("Error reading file1:", err)
            return
        }

        for a, b := range words1 {
            words1[a] = replaceQuotes(b)
        }

        /* ************************************************************************
            READ IN WORDS OF FILE 2
        ************************************************************************ */

        words2, err := readWords(file2)
        if err != nil {
            fmt.Println("Error reading file2:", err)
            return
        }

        for a, b := range words2 {
            words2[a] = replaceQuotes(b)
        }



        /* ************************************************************************
            SET UP AND RUN MAIN PROCESS
        ************************************************************************ */

        minLength := min(len(words1), len(words2))
        discrepancies := false

        fmt.Printf("Will begin editing %s according to %s\n\n", file1, file2)

        continueEditing := true

        var i int
        var j int

        for i,j = idx,jdx; i < minLength && continueEditing; {
            word1 := cleanWord(words1[i])
            word2 := cleanWord(words2[j])

            if word1 == "" {
                i++
                continue
            }
            if word2 == "" {
                j++
                continue
            }

            if word1 != word2 {
                discrepancies = true
                printSurroundingWords(words1, words2, i, j, file1, file2)

                // var choice int
                var choice string

                for {
                    // fmt.Println("Which version to keep? (1 for", file1, ", 2 for", file2, "):")
                    fmt.Printf("How to resolve? ('xy' numbers to advance cursor: x for file under edit, y for source file,  'a' to add missing token, 'e' to edit typo,  'd' to delete token, 'x' to delete current token from source, 'v' to save changes and quit)\n\n")
                    fmt.Scan(&choice)

                    // start := max(0, i-10)
                    // end := min(len(words1), i+11)
                    if choice == "s" {
                        //  will need to advance i to next equal words
                        i++
                        j++
                        break
                    } else if choice == "a" {
                        //  discrep is that file under edit is missing a token from the source
                        //  so add source token to the file under edit
                        words1 = append(words1[:i], append([]string{words2[j]}, words1[i:]...)...)
              
                        break
                    } else if choice == "e" {
                        //  discrep is just a mispelled word or the wrong word, but before and after, the text is good
                        //  so make words1[i] equal whatever is at words2[i]
                        words1[i] = words2[j]
                        //  resume comparison at i
                        break
                    } else if choice == "ex" {
                        //  e but apply the edit to the source file
                        words2[j] = words1[i]
                        //  resume comparison at i
                        break
                    } else if choice == "d" {
                        // surplus token in file under edit, delete that token
                        words1 = append(words1[0:i], words1[i+1:]...)

                        break
                    } else if choice == "v" {
                        // save current changes and exit
                        continueEditing = false
                        break
                    }  else if choice == "x" {
                        // delete token from source
                        words2 = append(words2[0:j], words2[j+1:]...)
                        break
                    } else {
                        digits, err := parseDigits(choice)
                        if err != nil {
                            fmt.Println("Error parsing digits:", err)
                            break
                        }
                        
                        i += digits[0]
                        j += digits[1]
                        break
                    }
                }
            } else {
                i++
                j++
            }
        }

        if discrepancies {
            if err := updateFile(file1, words1); err != nil {
                fmt.Println("Error updating file1:", err)
            }
            if err := updateFile(file2, words2); err != nil {
                fmt.Println("Error updating file2:", err)
            }
            fmt.Println("Files have been updated based on user choices.")
        } else {
            fmt.Println("Files are identical.")
        }

        fmt.Printf("\n\nleft at indexes [i = %d] [j = %d]\n\n", i, j)
    }

    func parseDigits(input string) ([]int, error) {
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

    func replaceQuotes(text string) string {
        // Replace straight double quotes with curly double quotes
        text = strings.ReplaceAll(text, "“", "\"")
        text = strings.ReplaceAll(text, "”", "\"")
    
        // Replace straight single quotes with curly single quotes
        text = strings.ReplaceAll(text, "‘", "'")
        text = strings.ReplaceAll(text, "’", "'")
    
        return text
    }


    func cleanWord(s string) string {
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
