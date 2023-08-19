package main

import (
	"strings"
	"testing"
)

// The tests in this suite can be heavily refactored, but this is a first attempt at testing.
// Just getting the vanilla feel of it.
// Also, most of the methods I will test now, I would normally not test because they are thin wrappers around standard library methods.

// Tests whether the structure of the csv file is read correctly.
func TestReadCSVFile(t *testing.T) {
	csvContent := "1+1,2\n2+3,5\n9+9,18"
	csvFile := strings.NewReader(csvContent)
	csvRows := ReadCSVFile(csvFile)

	if len(csvRows) != 3 {
		t.Fatalf("Expected rows count %d, actual count %d", 3, len(csvRows))
	}

	if csvRows[0][0] != "1+1" {
		t.Fatalf("Expected first question to be %s, but was %s", "1+1", csvRows[0][0])
	}

	if csvRows[0][1] != "2" {
		t.Fatalf("Expected first question's answer to be %s, but was %s", "2", csvRows[0][1])
	}

	if csvRows[1][0] != "2+3" {
		t.Fatalf("Expected second question to be %s, but was %s", "2+3", csvRows[1][0])
	}

	if csvRows[1][1] != "5" {
		t.Fatalf("Expected second question's answer to be %s, but was %s", "5", csvRows[1][1])
	}

	if csvRows[2][0] != "9+9" {
		t.Fatalf("Expected third question to be %s, but was %s", "9+9", csvRows[2][0])
	}

	if csvRows[2][1] != "18" {
		t.Fatalf("Expected third question's answer to be %s, but was %s", "18", csvRows[2][1])
	}
}

// Happy path, because we will validate the CSV before passing it to this function.
func TestConvertRowToProblem_CorrectSliceFormat(t *testing.T) {
	row := []string{"1+1", "2"}
	problem := ConvertRowToProblem(row)

	if problem.q != "1+1" || problem.a != "2" {
		t.Fatalf("Expected value for question is %s and for answer is %s, but was %s and %s",
			"1+1", "2", problem.q, problem.a)
	}
}

func TestValidateCSVRows_ValidCSVFormat(t *testing.T) {
	rows := [][]string{
		[]string{"1+1", "2"},
		[]string{"2+3", "5"},
	}

	if !ValidateCSVRows(rows) {
		t.Fatalf("CSV should have been considered valid.")
	}
}

func TestValidateCSVRows_InvalidCSVFormat(t *testing.T) {
	rows := [][]string{
		[]string{"1+1", "2"},
		[]string{"2+3"},
	}

	if ValidateCSVRows(rows) {
		t.Fatalf("CSV should have been considered invalid.")
	}
}
