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
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type problem struct {
	q string // question
	a string // answer
}

type flags struct {
	filePath string
	shuffle  bool
}

// Column count of a valid csv problem row record
const COLUMN_COUNT = 2

func printErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// Reads program's flags and populates a DS with their values.
// Flag messages and default values are set here.
func ReadProgramFlags() flags {
	filePath := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	shuffle := flag.String("shuffle", "0", "Whether to shuffle the questions. 1 for yes, 0 for no.")
	flag.Parse()

	shuffleBoolValue, err := strconv.ParseBool(*shuffle)

	if err != nil {
		printErrorAndExit("Could not parse boolean value for 'shuffle' flag")
	}

	return flags{
		filePath: *filePath,
		shuffle:  shuffleBoolValue,
	}
}

// Opens a file located at the specified path.
// If there is an error, on opening the file, a message is printed and the program exits.
// Thoughts: Usually, I won't write this method, because it is a thin wrapper around os.Open with no added value.
func OpenFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		printErrorAndExit(fmt.Sprintf("Failed to open the CSV file: %s\n", path))
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
		printErrorAndExit("Failed to parse the provided CSV file.")
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
		if len(row) != COLUMN_COUNT {
			return false
		}
	}

	return true
}

// A wrapper around convertLineToProblem for handling 2D slice to slice of problems mapping.
func ConvertRowsToProblems(rows [][]string) []problem {
	problems := make([]problem, len(rows))

	for i, row := range rows {
		problems[i] = ConvertRowToProblem(row)
	}

	return problems
}

// Compares the original answer to the one provided by the user.
// prob is the original question and answer
// r is the reader from which user input is read
// If there is a problem with reading user input or the answer does not match the actual answer, return false.
func CheckAnswer(prob *problem, r io.Reader) bool {
	var userAnswer string
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		userAnswer = scanner.Text()
	}

	if scanner.Err() != nil {
		fmt.Printf("Sorry, we couldn't understand your answer. We will consider it wrong, just in case :D")
		return false
	}

	return strings.TrimSpace(userAnswer) == prob.a
}

func GetProblemPromptMsg(problemNum int, question string) string {
	return fmt.Sprintf("Problem #%d: %s?\n", problemNum, question)
}

func GetResultMsg(correctAnswers int, questionCount int) string {
	return fmt.Sprintf("You got %d out of %d correctly!\n", correctAnswers, questionCount)
}

// This slice is passed by reference because we are shuffling it inplace and the result is visible from the outside scope.
func shuffleProblems(arr []problem) {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
}

func main() {
	flags := ReadProgramFlags()
	file := OpenFile(flags.filePath)
	csvRows := ReadCSVFile(file)

	if !ValidateCSVRows(csvRows) {
		printErrorAndExit("Invalid csv format")
	}

	problems := ConvertRowsToProblems(csvRows)

	if flags.shuffle {
		shuffleProblems(problems)
	}

	// Handle answers
	correctAnswers := 0
	for i, problem := range problems {
		fmt.Print(GetProblemPromptMsg(i+1, problem.q))
		if CheckAnswer(&problem, os.Stdin) {
			correctAnswers++
		}
	}

	fmt.Print(GetResultMsg(correctAnswers, len(problems)))
}
