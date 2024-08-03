package editingjob

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"poweredit/utils"
	"strconv"
	"strings"
)

var JOB_DIRECTORY string
var TEXT_DIRECTORY string

func init() {
	homedir, exists := os.LookupEnv("HOME")
	if !exists {
		log.Fatal("could not locate home directory by environment variable HOME")
	}
	
	JOB_DIRECTORY = homedir+"/.powerEdit/jobs"   //	use actual application file
	TEXT_DIRECTORY = homedir+"/.powerEdit/texts" //	use actual application file
}

func init() {
	createFileIfNotExist(JOB_DIRECTORY)
	createFileIfNotExist(TEXT_DIRECTORY)
}


type EditingJob struct { // TODO: have only a sinlge LatestEdition field
	name                 string
	editingFile           string // location of txt file
	sourceFile         string // location of txt file
    latestEditFile       string // location of edit file with latest edits
    latestSourceFile       string // location of srce file with latest edits
	latestEdition        int    // latest edition of txt file
	LastEditingIndex     int    // where to start processing txt file
	LastSourceIndex      int    // where to start processing txt file
}

func (ej *EditingJob) FieldNameSlice() []string {
    return []string{
        "name",
        "editing_file",
        "source_file",
        "latest_edit_file",
        "latest_source_file",
        "latest_edition",
        "last_editing_index",
        "last_source_index",
    }
}

func (ej *EditingJob) LatestEditFile() string {
	return ej.latestEditFile
}

func (ej *EditingJob) LatestSrceFile() string {
	return ej.latestSourceFile
}

func (ej *EditingJob) BumpEdition() {
	ej.latestEdition++
}

func (ej *EditingJob) ToStringSlice() []string {
    return []string{
        ej.name,
        ej.editingFile,
        ej.sourceFile,
        ej.latestEditFile,
        ej.latestSourceFile,
        fmt.Sprint(ej.latestEdition),
        fmt.Sprint(ej.LastEditingIndex),
        fmt.Sprint(ej.LastSourceIndex),
    }
}

func (ej *EditingJob) SaveLatestEdition(edits, source string) error {
	ej.latestEdition++
	newEdits := path.Join(TEXT_DIRECTORY,fmt.Sprintf("%d_%s",ej.latestEdition,path.Base(ej.editingFile)))
	newSource := path.Join(TEXT_DIRECTORY, fmt.Sprintf("%d_%s",ej.latestEdition,path.Base(ej.sourceFile)))
	if err := utils.UpdateFile(newEdits, edits); err != nil {
		ej.latestEdition--
		return fmt.Errorf("error updating %s: %v", newEdits, err)
	}
	if err := utils.UpdateFile(newSource,source); err != nil {
		ej.latestEdition--
		err := os.Remove(newEdits)
		return fmt.Errorf("error updating %s: %v", newSource, err)
	}
	return nil
}

func FromJobFile(jobfile string) (*EditingJob, error) {
	base := path.Base(jobfile)
	noext := strings.TrimRight(base, ".csv")
	return ReadEditingJob(path.Join(JOB_DIRECTORY,noext,jobfile))
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

	if headers[0] != "name" || headers[1] != "editing_file" || headers[2] != "source_file" || headers[3] != "latest_edit_file" ||
    headers[4] != "latest_source_file" || headers[5] != "latest_edition" || headers[6] != "last_editing_index" || headers[7] != "last_source_index" {
		return nil, fmt.Errorf("unexpected CSV headers for job file: %s\nconsider reviewing file", filename)
	}
	
	row := records[len(records)-1]

	latestEdition, err := strconv.Atoi(row[5])
	if err != nil {
		return nil, err
	}
	
	lastEditingIndex, err := strconv.Atoi(row[6])
	if err != nil {
		return nil, err
	}
	
	lastSourceIndex, err := strconv.Atoi(row[7])
	if err != nil {
		return nil, err
	}

	
	job := &EditingJob{
        name: row[0],
		editingFile:           row[1],
		sourceFile:         row[2],
        latestEditFile: row[3],
        latestSourceFile: row[4],
		latestEdition: latestEdition,
		LastEditingIndex:     lastEditingIndex,
		LastSourceIndex:      lastSourceIndex,
	}

	return job, nil
}

func UpdateEditingJob(filename string, job *EditingJob) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the job data
	record := job.ToStringSlice()
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func writeNewEditingJob(filename string, job *EditingJob) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("couldn't create %s: %v",filename,err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	headers := job.FieldNameSlice()
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	// Write the job data
	record := job.ToStringSlice()
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func FromEditAndSourceFiles(editFile, srceFile string) (*EditingJob, error) {
    baseEditName := path.Base(editFile)
    baseSrceName := path.Base(srceFile)
    shortEditFileName := strings.TrimRight(baseEditName, path.Ext(editFile))
    shortSrceFileName := strings.TrimRight(baseSrceName, path.Ext(srceFile))
    jobname := strings.TrimSpace(fmt.Sprintf("edit_%s_by_%s", shortEditFileName, shortSrceFileName))

    newJob := EditingJob{
        name:                 jobname,
        editingFile:           editFile,
        sourceFile:         srceFile,
        latestEditFile:       path.Join(TEXT_DIRECTORY, "0_"+baseEditName),
        latestSourceFile:       path.Join(TEXT_DIRECTORY, "0_"+baseSrceName),
        latestEdition:        0,
        LastEditingIndex:     0,
        LastSourceIndex:      0,
    }

    err := writeAllJobFiles(&newJob)
    if err != nil {
        return nil,
        fmt.Errorf("couldn't create new editing job from %s and %s: %v", editFile, srceFile, err)
    }

    return &newJob, nil
}

func writeAllJobFiles(job *EditingJob) error {

	//  create csv job filename eg. ~/.powerEdit/jobs/edit_badfoo_by_goodfoo/edit_badfoo_by_goodfoo.csv
	newJobfileName := path.Join(JOB_DIRECTORY, job.name, job.name+".csv")

	createFileIfNotExist(path.Join(JOB_DIRECTORY,job.name))

	//  write job csv
	if err := writeNewEditingJob(newJobfileName, job); err != nil {
		return fmt.Errorf("failed to create new csv jobfile: %v", err)
	}

	//  read the original edit and source files
	editFileContent, err := os.ReadFile(job.editingFile)
	if err != nil {
		return fmt.Errorf("couldn't read EditingFile while creating new edit job: %v", err)
	}

	sourceFileContent, err := os.ReadFile(job.sourceFile)
	if err != nil {
		return fmt.Errorf("couldn't read SourceFile while creating new edit job: %v", err)
	}

	//  write the first editing version file for the file to edit
    //      will be used for first edit session, wherefrom edits will be saved to v_1. So v_0 also serves as backup for originals
	if err := os.WriteFile(job.latestEditFile, editFileContent, 0644); err != nil {
		return fmt.Errorf("couldn't write version_0 edit file: %v", err)
	}

	//  write the first editing version file for the source of edits
    //      will be used for first edit session, wherefrom edits will be saved to v_1. So v_0 also serves as backup for originals
	if err := os.WriteFile(job.latestSourceFile, sourceFileContent, 0644); err != nil {
		return fmt.Errorf("couldn't write version_0 source file: %v", err)
	}

	return nil
}

func getAllJobs() ([]fs.DirEntry, error) {
	files, err := os.ReadDir(JOB_DIRECTORY)
	if err != nil {
		return nil, fmt.Errorf("could not get jobs: %v", err)
	}

	return files, nil
}

func JobExists(jobname string) (bool, error) {
	jobs, err := getAllJobs()
	if err != nil {
		return false, fmt.Errorf("could not display jobs: %v", err)
	}

	for _, job := range jobs {
		if job.IsDir() {
			if jobname == job.Name() {
				return true, nil
			}
		}
	}

	return false, nil
}

func DisplayJobs() error {
	files, err := getAllJobs()
	if err != nil {
		return fmt.Errorf("could not display jobs: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name())
		}
	}

	return nil
}

func createFileIfNotExist(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(filename, os.ModePerm)
		if err != nil {
			return fmt.Errorf("%s doesnt exist but could not be created: %v", filename, err)
		}
	}
	return nil
}