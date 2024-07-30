package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

type EditingJob struct {        // TODO: have only a sinlge LatestEdition field
	EditingFile          string // location of txt file
	SourceFile           string // location of txt file
	LatestEditingEdition int    // latest edition of txt file
	LatestSourceEdition  int    // latest edition of txt file
	LastEditingIndex     int    // where to start processing txt file
	LastSourceIndex      int    // where to start processing txt file
}

func (ej *EditingJob) FileToEdit() string {
    filedir := path.Dir(ej.EditingFile)
    basename := path.Base(ej.EditingFile)
    return fmt.Sprintf("%s/%d_%s",filedir,ej.LatestEditingEdition,basename)
}

func (ej *EditingJob) FileOfSource() string {
    filedir := path.Dir(ej.SourceFile)
    basename := path.Base(ej.SourceFile)
    return fmt.Sprintf("%s/%d_%s",filedir,ej.LatestSourceEdition,basename)
}

func (ej *EditingJob) BumpEdition() {
    ej.LatestEditingEdition++
    ej.LatestSourceEdition++
}

func ReadEditingJob(filename string) (*EditingJob, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(records) < 2 {
        return nil, fmt.Errorf("CSV file must have at least one record")
    }
	
	headers := records[0]

	if headers[0] != "editing_file" || headers[1] != "source_file" || headers[2] != "latest_editing_edition" ||
	headers[3] != "latest_source_edition" || headers[4] != "last_editing_index" || headers[5] != "last_source_index" {
		return nil, fmt.Errorf("unexpected CSV headers for job file: %s\nconsider reviewing file",filename)
	}

    row := records[1]
    latestEditingEdition, err := strconv.Atoi(row[2])
    if err != nil {
        return nil, err
    }
    latestSourceEdition, err := strconv.Atoi(row[3])
    if err != nil {
        return nil, err
    }
    lastEditingIndex, err := strconv.Atoi(row[4])
    if err != nil {
        return nil, err
    }
    lastSourceIndex, err := strconv.Atoi(row[5])
    if err != nil {
        return nil, err
    }

    job := &EditingJob{
        EditingFile:              row[0],
        SourceFile:               row[1],
        LatestEditingEdition: latestEditingEdition,
        LatestSourceEdition:  latestSourceEdition,
        LastEditingIndex:     lastEditingIndex,
        LastSourceIndex:      lastSourceIndex,
    }

    return job, nil
}

func UpdateEditingJob(filename string, job *EditingJob) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }

    if len(records) < 2 {
        return fmt.Errorf("CSV file must have at least one record")
    }

    records[1][0] = job.EditingFile
    records[1][1] = job.SourceFile
    records[1][2] = strconv.Itoa(job.LatestEditingEdition)
    records[1][3] = strconv.Itoa(job.LatestSourceEdition)
    records[1][4] = strconv.Itoa(job.LastEditingIndex)
    records[1][5] = strconv.Itoa(job.LastSourceIndex)

    file, err = os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, record := range records {
        err := writer.Write(record)
        if err != nil {
            return err
        }
    }

    return nil
}

func WriteNewEditingJob(filename string, job *EditingJob) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write the header
    header := []string{"editing", "source", "latest_editing_edition", "latest_source_edition", "last_editing_index", "last_source_index"}
    err = writer.Write(header)
    if err != nil {
        return err
    }

    // Write the job data
    record := []string{
        job.EditingFile,
        job.SourceFile,
        strconv.Itoa(job.LatestEditingEdition),
        strconv.Itoa(job.LatestSourceEdition),
        strconv.Itoa(job.LastEditingIndex),
        strconv.Itoa(job.LastSourceIndex),
    }
    err = writer.Write(record)
    if err != nil {
        return err
    }

    return nil
}

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
func UpdateFile(filename string, content string) error {
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
