package editingjob

import (
	"fmt"
	"os"
	"path"
	"slices"
	"testing"
)

var (
	test_name = "edit_gutenberg-iliad_by_ia-iliad"
	test_editingFile = "test/testpowereditdir/testtexts/gutenberg-iliad.txt"
	test_sourceFile = "test/testpowereditdir/testtexts/ia-iliad.txt"
	test_latestEditFile = "test/testpowereditdir/testtexts/5_gutenberg-iliad.txt"
	test_latestSourceFile = "test/testpowereditdir/testtexts/5_ia-iliad.txt"
	test_latestEdition = 5
	test_LastEditingIndex = 24515
	test_LastSourceIndex = 24408

	test_newjob_name = "edit_gutn_by_iarc"
	test_newjob_editingFile = "test/testpowereditdir/testtexts/gutn.txt"
	test_newjob_sourceFile = "test/testpowereditdir/testtexts/iarc.txt"
	test_newjob_latestEditFile = "test/testpowereditdir/testtexts/0_gutn.txt"
	test_newjob_latestSourceFile = "test/testpowereditdir/testtexts/0_iarc.txt"
	test_newjob_latestEdition = 0
	test_newjob_lastEditingIndex = 0
	test_newjob_lastSourceIndex = 0


	mockExistingEditingJob = EditingJob {
		name: test_name,
		editingFile: test_editingFile,
		sourceFile: test_sourceFile,
		latestEditFile: test_latestEditFile,
		latestSourceFile: test_latestSourceFile,
		latestEdition: test_latestEdition,
		LastEditingIndex: test_LastEditingIndex,
		LastSourceIndex: test_LastSourceIndex,
	}

	mockNewEditingJob = EditingJob {
		name: test_newjob_name,
		editingFile: test_newjob_editingFile,
		sourceFile: test_newjob_sourceFile,
		latestEditFile: test_newjob_latestEditFile,
		latestSourceFile: test_newjob_latestSourceFile,
		latestEdition: test_newjob_latestEdition,
		LastEditingIndex: test_newjob_lastEditingIndex,
		LastSourceIndex: test_newjob_lastSourceIndex,
	}

	TEST_JOB_DIRECTORY = "test/testpowereditdir/testjobs"
	TEST_TEXT_DIRECTORY = "test/testpowereditdir/testtexts"
	TEST_JOBFILE = "edit_gutenberg-iliad_by_ia-iliad.csv"

	TEST_EDIT_FILE_BASE_NOEXT = "gutenberg-iliad"
	TEST_SOURCE_FILE_BASE_NOEXT = "ia-iliad"
	TEST_EDIT_FILE_BASE = TEST_EDIT_FILE_BASE_NOEXT+".txt"
	TEST_SOURCE_FILE_BASE = TEST_SOURCE_FILE_BASE_NOEXT+".txt"

	TEST_NEWJOB_EDIT_FILE_BASE_NOEXT = "gutn"
	TEST_NEWJOB_SOURCE_FILE_BASE_NOEXT = "iarc"
	TEST_NEWJOB_EDIT_FILE_BASE = TEST_NEWJOB_EDIT_FILE_BASE_NOEXT+".txt"
	TEST_NEWJOB_SOURCE_FILE_BASE = TEST_NEWJOB_SOURCE_FILE_BASE_NOEXT+".txt"
)

func TestToStringSlice(t *testing.T) {
	res := mockExistingEditingJob.ToStringSlice()
	expect := []string{
		test_name,
		test_editingFile,
		test_sourceFile,
		test_latestEditFile,
		test_latestSourceFile,
		fmt.Sprint(test_latestEdition),
		fmt.Sprint(test_LastEditingIndex),
		fmt.Sprint(test_LastSourceIndex),
	}

	if !slices.Equal[[]string, string](res, expect) {
		t.Errorf("got: %v, want %v", res, expect)
	}

	notExpected := []string{
		test_name,
		test_editingFile,
		test_sourceFile,
		test_latestEditFile,
		test_latestSourceFile,
		fmt.Sprint(test_latestEdition),
		fmt.Sprint(test_LastEditingIndex),
		"3",
	}

	if slices.Equal[[]string, string](res, notExpected) {
		t.Error("test should not have passed")
	}
}

func TestFieldNameSlice(t *testing.T) {
	res := mockExistingEditingJob.FieldNameSlice()
	expect := []string{
        "name",
        "editing_file",
        "source_file",
        "latest_edit_file",
        "latest_source_file",
        "latest_edition",
        "last_editing_index",
        "last_source_index",
    }

	if !slices.Equal[[]string, string](res, expect) {
		t.Errorf("got: %v, want: %v", res, expect)
	}

	notExpected := []string{
        "name",
        "editing_file",
        "source_file",
        "latest_edit_file",
        "latest_source_file",
        "latest_editing_edition",
        "last_editing_index",
        "last_source_index",
    }

	if slices.Equal[[]string, string](res, notExpected) {
		t.Error("test should not have passed")
	}
}

func TestLatestEditFile(t *testing.T) {
	res := mockExistingEditingJob.LatestEditFile()
	expect := test_latestEditFile

	if res != expect {
		t.Errorf("got: %s, want: %s", res, expect)
	}
}

func TestLatestSrceFile(t *testing.T) {
	res := mockExistingEditingJob.LatestSrceFile()
	expect := test_latestSourceFile

	if res != expect {
		t.Errorf("got: %s, want: %s", res, expect)
	}
}

func TestBumpEdition(t *testing.T) {
	expect := test_latestEdition + 1
	mockExistingEditingJob.BumpEdition()
	res := mockExistingEditingJob.latestEdition

	defer func() {
		mockExistingEditingJob.latestEdition = 5
	}()

	if res != expect {
		t.Errorf("got: %d, want: %d", res, expect)
	}
}

func TestFromJobFile(t *testing.T) {
	JOB_DIRECTORY = TEST_JOB_DIRECTORY
	TEXT_DIRECTORY = TEST_TEXT_DIRECTORY

	res, err := FromJobFile(TEST_JOBFILE)
	if err != nil {
		t.Errorf("test resulted in error: %v", err)
	}

	if *res != mockExistingEditingJob {
		t.Errorf("\ngot:  %#v, \nwant: %#v\n", *res, mockExistingEditingJob)
	}
}

func TestFromEditAndSourceFiles(t *testing.T) {
	JOB_DIRECTORY = TEST_JOB_DIRECTORY
	TEXT_DIRECTORY = TEST_TEXT_DIRECTORY

	newJobDir := path.Join(JOB_DIRECTORY, test_newjob_name)
	newJobFilePath := path.Join(newJobDir, test_newjob_name+".csv")

	deleteFile(newJobFilePath)
	deleteFile(newJobDir)

	fullPathEditFile := path.Join(TEST_TEXT_DIRECTORY,TEST_NEWJOB_EDIT_FILE_BASE)
	fullPathSourceFile := path.Join(TEST_TEXT_DIRECTORY,TEST_NEWJOB_SOURCE_FILE_BASE)


	res, err := FromEditAndSourceFiles(fullPathEditFile, fullPathSourceFile)
	if err != nil {
		t.Errorf("test resulted in error: %v", err)
	}

	if !deleteFile(newJobFilePath) {
		t.Errorf("new jobfile %s was not created", newJobFilePath)
	}

	if !deleteFile(newJobDir) {
		t.Errorf("new job dir %s was not created", newJobDir)
	}

	if *res != mockNewEditingJob {
		t.Errorf("\ngot:  %#v\n, \nwant: %#v\n", *res, mockNewEditingJob)
	}
}

func deleteFile(filename string) bool {
    err := os.Remove(filename)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil
}