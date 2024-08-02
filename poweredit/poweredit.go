package poweredit


import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"poweredit/editingjob"
	"poweredit/textwords"
	"poweredit/utils"
	"strings"
)

var jobfile string
var editIndexFlag int
var sourceIndexFlag int

var jobdata *editingjob.EditingJob

func init() {
	flag.IntVar(&editIndexFlag, "ei", -1, "editing index (ei) - location to start edit comparison in editing file")
	flag.IntVar(&sourceIndexFlag, "si", -1, "source index (si) - location to start edit comparison in source file")
}

func initJob() {
	args := flag.Args()
	argln := len(args)

	if argln == 1 {

		if args[0] == "jobs" {
			utils.ClearScreen()
			fmt.Printf("All available jobs:\n\n")
			err := editingjob.DisplayJobs()
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(0)
		}
	
		if strings.HasSuffix(args[0], ".csv") {
			jobfile = args[0]
			job, err := editingjob.FromJobFile(jobfile)
			if err != nil {
				fmt.Printf("Couldn't locate jobfile %s\n%v", jobfile, err)
				os.Exit(0)
			}
			jobdata = job
			return
		}

		validJob, _ := editingjob.JobExists(args[0])

		if (validJob) {
			jobfile = args[0]+".csv"
			job, err := editingjob.FromJobFile(jobfile)
			if err != nil {
				fmt.Printf("Couldn't locate jobfile %s\n%v", jobfile, err)
				os.Exit(0)
			}
			jobdata = job
			return
		}


	} else if argln == 2 {
		newEditingFile, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Printf("tried to create new job with %s editing file but could not find absolute path to the file: %v", args[0], err)
		}
		newSourceFile, err  := filepath.Abs(args[1])
		if err != nil {
			fmt.Printf("tried to create new job with %s source file but could not find absolute path to the file: %v", args[1], err)
		}
		
		job, err := editingjob.FromEditAndSourceFiles(newEditingFile, newSourceFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		jobdata = job
		return
	}

	fmt.Println("invalid use. Must provide jobfile OR new-editing-file AND new-source-file")
	os.Exit(0)

}

func Run(args []string) {

	

	/* ************************************************************************
		PARSE COMMAND LINE FLAGS
	************************************************************************ */
	flag.Parse()

	

	/* ************************************************************************
		EVALUATE COMMAND LINE ARGS TO CREATE EDITING JOB STRUCT
	************************************************************************ */
	initJob()

	

	/* ************************************************************************
		SET STARTING INDEXES FOR EDITING JOB
	************************************************************************ */
	var i int
	var j int

	if editIndexFlag >= 0 {
		i = editIndexFlag
	} else {
		i = jobdata.LastEditingIndex
	}

	if sourceIndexFlag >= 0 {
		j = sourceIndexFlag
	} else {
		j = jobdata.LastSourceIndex
	}

	/* ************************************************************************
		CREATE TEXTWORDS STRUCTS FOR EDIT AND SOURCE FILES
	************************************************************************ */
		editWords, err := textwords.FromFile(jobdata.LatestEditFile())
		if err != nil {
			fmt.Println("Error getting edit words:", err)
			return
		}

		sourceWords, err := textwords.FromFile(jobdata.LatestSrceFile())
		if err != nil {
			fmt.Println("Error getting source words:", err)
			return
		}

	/* ************************************************************************
		ITEREATE OVER TEXTWORDS UNTIL END OF JOB
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
			printDisplay(jobdata,editWords,sourceWords,i,j)
			
			var choice string

			for {
				printResolutionOptions()
				fmt.Scan(&choice)

				if choice == "a" {
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
				} else if choice == "q" {
					// quit without saving any edits
					utils.ClearScreen()
					os.Exit(0)
				} else if choice == "v" {
					// save current changes and exit
					continueEditing = false
					break
				} else if choice == "x" {
					// delete token from source
					sourceWords.Delete(j)
					break
				} else if choice == "me" {
					var customWord string
					var confirm string
					// manually enter word and edit both by this word
					for {
						fmt.Print("enter word to edit both by: ")
						fmt.Scan(&customWord)
						fmt.Printf("save '%s' to both indexes? (y/n): ", customWord)
						fmt.Scan(&confirm)
						if strings.ToLower(confirm) == "y" {
							break
						}
						printDisplay(jobdata,editWords,sourceWords,i,j)
					}

					editWordLoc.W = customWord
					sourceWordLoc.W = customWord
					editWords.Edit(i, editWordLoc)
					sourceWords.Edit(j, sourceWordLoc)
					break
				} else {
					digits, err := utils.ParseDigits(choice)
					if err != nil {
						fmt.Printf("Not a valid command: %s", choice)
						break
					}

					if len(digits) == 1 {
						i += digits[0]
						break
					} else {
						i += digits[0]
						j += digits[1]
						break
					}
				}
			}
		} else {
			i++
			j++
		}
	}

	if discrepancies {
		jobdata.BumpEdition()
		if err := utils.UpdateFile(jobdata.LatestEditFile(), editWords.Text()); err != nil {
			fmt.Printf("Error updating %s: %v", jobdata.LatestEditFile(), err)
		}
		if err := utils.UpdateFile(jobdata.LatestSrceFile(), sourceWords.Text()); err != nil {
			fmt.Printf("Error updating %s: %v", jobdata.LatestSrceFile(), err)
		}

		jobdata.LastEditingIndex = i
		jobdata.LastSourceIndex = j

		editingjob.UpdateEditingJob(jobfile, jobdata)

		fmt.Println("Files have been updated based on user choices.")

	} else {
		fmt.Println("Files are identical.")
	}

	fmt.Printf("\n\nleft at indexes [i = %d] [j = %d]\n\n", i, j)
}



func printDisplay(job *editingjob.EditingJob, ew, sw *textwords.TextWords, i, j int) {
	utils.Display(fmt.Sprintf("\n\tediting %s  by  %s\n\n\n", path.Base(job.LatestEditFile()), path.Base(job.LatestSrceFile())))
	fmt.Printf("\tDISCREPANCY:\n\n\tfile under edit: %s\n\tsource file:     %s\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n", ew.SurroundingText(i, 10), sw.SurroundingText(j, 10))
}

func printResolutionOptions() {
	fmt.Printf(
		"\tHow to resolve?\n" +
			"\t<xy|x> - enter two numbers to advance cursors: x for file under edit, y for source file\n" +
			"\t\ta single digit entry will advance cursor for file under edit by x\n" +
			"\ta - to add missing token to file under edit\n" +
			"\te - edit typo, sets current word of file under edit to current word of source file\n" +
			"\tex - edit typo in source, sets current word of source file to current word of file under edit\n" +
			"\tme - manually enter a custom word set current token for file under edit and source file to this word\n" +
			"\td - delete token from file under edit\n" +
			"\tx - delete current token from source file\n" +
			"\tv - save changes and quit\n" +
			"\tq - quit without saving any changes made\n\n\tenter selection: ")
}