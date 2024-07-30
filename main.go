package main

import (
	"compare/textwords"
	"compare/utils"
	"fmt"
	"os"
	"strconv"
)

func main() {

	/* ************************************************************************
		GET OPTIONS SET
	************************************************************************ */
	if len(os.Args) < 3 {
		fmt.Println("Usage: compare <file to edit> <'source file' upon which edits are based>")
		return
	}

	i := 0
	j := 0

	file1 := os.Args[1] // file to edit
	file2 := os.Args[2] // source file

	i, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Couldn't parse idx, using 0")
		i = 0
	}

	j, err = strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Couldn't parse jdx, using 0")
		j = 0
	}

	/* ************************************************************************
		READ IN WORDS OF FILE 1
	************************************************************************ */

	editWords, err := textwords.FromFile(file1)
	if err != nil {
		fmt.Println("Error getting edit words:", err)
		return
	}

	/* ************************************************************************
		READ IN WORDS OF FILE 2
	************************************************************************ */

	sourceWords, err := textwords.FromFile(file2)
	if err != nil {
		fmt.Println("Error getting source words:", err)
		return
	}

	/* ************************************************************************
		SET UP AND RUN MAIN PROCESS
	************************************************************************ */

	minLength := min(editWords.Len(), sourceWords.Len())
	
	discrepancies := false
	continueEditing := true

	for ; i < minLength && continueEditing; {
		editWordLoc 	:= editWords.GetWord(i)
		sourceWordLoc 	:= sourceWords.GetWord(j)

		editWord := utils.CleanWord(editWordLoc.W)
		sourceWord := utils.CleanWord(sourceWordLoc.W)

		if editWord == "" {
			i++
			continue
		}
		if sourceWord == "" {
			j++
			continue
		}

		if editWord != sourceWord {
			discrepancies = true
			fmt.Printf("discrepancy:\ntxt: %s\nsrc: %s\n",editWords.SurroundingText(i,10),sourceWords.SurroundingText(j,10))

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
					editWords.Insert(sourceWordLoc,i)

					break
				} else if choice == "e" {
					//  discrep is just a mispelled word or the wrong word, but before and after, the text is good
					//  so make words1[i] equal whatever is at words2[i]
					editWords.Edit(i, sourceWordLoc)
					//  resume comparison at i
					break
				} else if choice == "ex" {
					//  e but apply the edit to the source file
					sourceWords.Edit(j, editWordLoc)
					//  resume comparison at i
					break
				} else if choice == "d" {
					// surplus token in file under edit, delete that token
					editWords.Delete(i)

					break
				} else if choice == "v" {
					// save current changes and exit
					continueEditing = false
					break
				} else if choice == "x" {
					// delete token from source
					sourceWords.Delete(j)
					break
				} else {
					digits, err := utils.ParseDigits(choice)
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
		if err := utils.UpdateFile("mod_"+file1, editWords.Text()); err != nil {
			fmt.Println("Error updating file1:", err)
		}
		if err := utils.UpdateFile("mod_"+file2, sourceWords.Text()); err != nil {
			fmt.Println("Error updating file2:", err)
		}
		fmt.Println("Files have been updated based on user choices.")
	} else {
		fmt.Println("Files are identical.")
	}

	fmt.Printf("\n\nleft at indexes [i = %d] [j = %d]\n\n", i, j)
}