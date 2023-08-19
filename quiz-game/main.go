// Implements the quiz game entirely.
// For simplicity, all functions are exported, so they can be easily tested.
// A more correct approach would be to have a separate package for related functionality and have tests for it.
// But here we will test everything together.
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

type problem struct {
	q string
	a string
}

type flags struct {
	filePath string
}

// Reads program's flags and populates a DS with their values.
// Flag messages and default values are set here.
func ReadProgramFlags() flags {
	filePath := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	flag.Parse()

	return flags{
		filePath: *filePath,
	}
}

// Opens a file located at the specified path.
// If there is an error, on opening the file, a message is printed and the program exits.
func OpenFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s\n", path)
		os.Exit(1)
	}

	return file
}

// Reads a csv formatted buffer into a 2D slice of strings.
// The rows in the slice represent the rows in the csv and the columns of the slice represent the columns of the csv.
// If the csv file is malformed a message is printed and the program exits.
func ReadCSVFile(file io.Reader) [][]string {
	csvReader := csv.NewReader(file)

	lines, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("Failed to parse the provided CSV file.")
		os.Exit(1)
	}

	return lines
}

// Does a mapping between a slice's elements and a problem.
// The slice's elements are an ordered set of the problem's fields as they appear in the problem struct.
func ConvertRowToProblem(row []string) problem {
	return problem{
		q: row[0],
		a: row[1],
	}
}

// If the row has two columns, we consider it valid.
func ValidateCSVRows(rows [][]string) bool {
	for _, row := range rows {
		if len(row) != 2 {
			return false
		}
	}

	return true
}

// A wrapper around convertLineToProblem for handling 2D slice to slice of problems mapping.
func ConvertRowsToProblems(rows [][]string) []problem {
	problems := make([]problem, len(rows))

	for i, line := range rows {
		problems[i] = ConvertRowToProblem(line)
	}

	return problems
}

// Compares the original answer to the one provided by the user.
// probNum is the display number of the problem
// prob is the original question and answer
// r is the reader from which user input is read
// If there is a problem with reading user input or the answer does not match the actual answer, return false.
func CheckAnswer(probNum int, prob *problem, r io.Reader) bool {
	var userAnswer string
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		userAnswer = scanner.Text()
	}

	if scanner.Err() != nil {
		fmt.Printf("Sorry, we couldn't understand your answer")
		return false
	}

	return userAnswer == prob.a
}

func GetProblemPrompt(problemNum int, question string) string {
	return fmt.Sprintf("Problem #%d: %s?\n", problemNum, question)
}

func GetResultMessage(correctAnswers int, questionCount int) string {
	return fmt.Sprintf("You got %d out of %d correctly!\n", correctAnswers, questionCount)
}

func main() {
	flags := ReadProgramFlags()
	file := OpenFile(flags.filePath)
	csvRows := ReadCSVFile(file)

	if ValidateCSVRows(csvRows) {
		fmt.Println("Invalid csv format")
		os.Exit(1)
	}

	problems := ConvertRowsToProblems(csvRows)

	// Handle answers
	correctAnswers := 0
	for i, problem := range problems {
		problemNum := i + 1
		fmt.Print(GetProblemPrompt(problemNum, problem.q))
		if CheckAnswer(problemNum, &problem, os.Stdin) {
			correctAnswers++
		}
	}

	fmt.Print(GetResultMessage(correctAnswers, len(problems)))
}
