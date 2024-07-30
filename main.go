package main

import (
	"compare/textwords"
	"compare/utils"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var jobfile string
var newEditingFile string
var newSourceFile string
var editingFileStartingIndex int
var sourceFileStartingIndex int

var jobdata *utils.EditingJob

func init() {
	flag.StringVar(&jobfile, "jf", "", "jobfile (jf) - location of csv file which holds data for editing job")
	flag.StringVar(&newEditingFile, "ef", "", "editing file (ef) - for initializing new jobs; location of editing file; only valid if new-source-file is provided as well; will override jobfile")
	flag.StringVar(&newSourceFile, "sf", "", "source file (sf) - for initializing new jobs; location of source file; only valid if new-editing-file is provided as well; will override jobfile")
	flag.IntVar(&editingFileStartingIndex, "ei", -1, "editing index (ei) - location to start edit comparison in editing file")
	flag.IntVar(&sourceFileStartingIndex, "si", -1, "source index (si) - location to start edit comparison in source file")
}

func initJob() {
	if newEditingFile != "" && newSourceFile != "" {
		newJob := utils.EditingJob{
			EditingFile: newEditingFile,
			SourceFile: newSourceFile,
			LatestEditingEdition: 0,
			LatestSourceEdition: 0,
			LastEditingIndex: 0,
			LastSourceIndex: 0,
		}
		editf := strings.TrimRight(path.Base(newEditingFile),path.Ext(newEditingFile))
		srcef := strings.TrimRight(path.Base(newSourceFile),path.Ext(newSourceFile))
		jobname := fmt.Sprintf("edit_%s_by_%s.csv", editf, srcef)

		err := utils.WriteNewEditingJob(jobname, &newJob)
		if err != nil {
			fmt.Printf("Couldn't create new jobfile %s\n%v", jobname,err)
			os.Exit(0)
		}

		jobdata = &newJob
	} else if jobfile != "" {
		job, err := utils.ReadEditingJob(jobfile)
		if err != nil {
			fmt.Printf("Couldn't locat jobfile %s\n%v", jobfile,err)
			os.Exit(0)
		}
		jobdata = job
	} else {
		fmt.Println("invalid use. Must provide jobfile OR new-editing-file AND new-source-file")
		os.Exit(0)
	}
}

func main() {

	flag.Parse()

	initJob()

	var i int
	var j int

	if editingFileStartingIndex >= 0 {
		i = editingFileStartingIndex
	} else {
		i = jobdata.LastEditingIndex
	}

	if sourceFileStartingIndex >= 0 {
		j = sourceFileStartingIndex
	} else {
		j = jobdata.LastSourceIndex
	}

	/* ************************************************************************
		READ IN WORDS OF FILE 1
	************************************************************************ */

	editWords, err := textwords.FromFile(jobdata.FileToEdit())
	if err != nil {
		fmt.Println("Error getting edit words:", err)
		return
	}

	/* ************************************************************************
		READ IN WORDS OF FILE 2
	************************************************************************ */

	sourceWords, err := textwords.FromFile(jobdata.FileOfSource())
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

	for i < minLength && continueEditing {
		editWordLoc := editWords.GetWord(i)
		sourceWordLoc := sourceWords.GetWord(j)

		editWord := utils.ReplaceQuotes(editWordLoc.W)
		sourceWord := utils.ReplaceQuotes(sourceWordLoc.W)

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
			fmt.Printf("discrepancy:\ntxt: %s\nsrc: %s\n", editWords.SurroundingText(i, 10), sourceWords.SurroundingText(j, 10))

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
					editWords.Insert(sourceWordLoc, i)

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
		jobdata.BumpEdition()
		if err := utils.UpdateFile(jobdata.FileToEdit(), editWords.Text()); err != nil {
			fmt.Printf("Error updating %s: %v",jobdata.FileToEdit(), err)
		}
		if err := utils.UpdateFile(jobdata.FileOfSource(), sourceWords.Text()); err != nil {
			fmt.Println("Error updating %s: %v",jobdata.FileOfSource(), err)
		}

		jobdata.LastEditingIndex = i
		jobdata.LastSourceIndex = j

		utils.UpdateEditingJob(jobfile, jobdata)

		fmt.Println("Files have been updated based on user choices.")


	} else {
		fmt.Println("Files are identical.")
	}

	fmt.Printf("\n\nleft at indexes [i = %d] [j = %d]\n\n", i, j)
}
