// Implements the quiz game entirely.
// For simplicity, all functions are exported, so they can be easily tested.
// It is a small program for the purposes of learning Go's idioms, syntax and standard library.
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	flags := readProgramFlags()
	file := openFile(flags.csvPath)
	csvRows := ReadCSV(file)

	if !ValidateColumnCount(csvRows, 2) {
		exit("Invalid csv format.")
	}

	problems := ConvertRowsToProblems(csvRows)

	if flags.shuffle {
		shuffle(problems)
	}

	timer := time.NewTimer(time.Duration(flags.time) * time.Second)

	// Handle user answers
	correct := 0
problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s?\n", i+1, problem.question)
		answer := make(chan string)

		// get input async and await on a channel
		go func() {
			var txt string
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				txt = scanner.Text()
			}
			answer <- txt
		}()
		select {
		case <-timer.C: // Time is up!
			fmt.Println()
			break problemloop
		case answer := <-answer: // The user answered!
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You got %d out of %d correctly!\n", correct, len(problems))
}

type Problem struct {
	question string
	answer   string
}

type quizFlags struct {
	csvPath string
	shuffle bool
	// Time limit for the whole quiz
	time int
}

func exitErr(msg string, err error) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s. Original error: %s", msg, err.Error()))
	os.Exit(1)
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// Reads program's flags and populates an object.
// CLI flag messages and default values are set here.
func readProgramFlags() quizFlags {
	csvPath := flag.String("csv", "problems.csv", "Path to a csv file in the format 'question,answer'.")
	shuffle := flag.Bool("shuffle", false, "Whether to shuffle the questions.")
	time := flag.Int("timeLimit", 30, "In Seconds, for the whole quiz.")
	flag.Parse()

	return quizFlags{
		csvPath: *csvPath,
		shuffle: *shuffle,
		time:    *time,
	}
}

// Opens a file located at the specified path.
// If there is an error, on opening the file, a message is printed and the program exits.
// Thoughts: Usually, I won't write this method, because it is a thin wrapper around os.Open with no added value.
func openFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		exitErr(fmt.Sprintf("Failed to open the CSV file: %s\n", path), err)
	}

	return file
}

// Reads a csv formatted buffer into a 2D slice of strings.
// The rows in the slice represent the rows in the csv and the columns of the slice represent the columns of the csv.
// If the csv data is malformed, a message is printed and the program exits.
// Calling this method creates an extra copy of the 2D slice, but the questions won't be that many.
// This method improves readability, which is more important than the supposed negative performance impact.
func ReadCSV(r io.Reader) [][]string {
	csv := csv.NewReader(r)
	rows, err := csv.ReadAll()

	if err != nil {
		exitErr("Failed to read CSV data.", err)
	}

	return rows
}

// Does a mapping between a slice's elements and a problem.
// The slice's elements are an ordered set of the problem's fields as they appear in the problem struct.
func MapFields(fields []string) Problem {
	return Problem{
		question: fields[0],
		answer:   fields[1],
	}
}

// If the row has two columns, we consider it valid.
func ValidateColumnCount(rows [][]string, count int) bool {
	for _, r := range rows {
		if len(r) != count {
			return false
		}
	}

	return true
}

// Each row is mapped to a problem object
func ConvertRowsToProblems(rows [][]string) []Problem {
	p := make([]Problem, len(rows))

	for i, r := range rows {
		p[i] = MapFields(r)
	}

	return p
}

// This slice is passed by reference because we are shuffling it inplace and the result is visible from the outside scope.
func shuffle(p []Problem) {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
}
