package main

import (
	"reflect"
	"strings"
	"testing"
)

// Tests whether the structure of the csv file is read correctly.
// This tests the csv package itself so it does not make sense to have this test.
// This is just explorational testing.
func TestReadCSV(t *testing.T) {
	tests := map[string]struct {
		input string
		want  [][]string
		err   bool
	}{
		"single line csv": {
			input: "1+1,2",
			want:  [][]string{[]string{"1+1", "2"}},
		},
		"multiline csv": {
			input: "1+1,2\n2+3,5\n9+9,18",
			want:  [][]string{[]string{"1+1", "2"}, []string{"2+3", "5"}, []string{"9+9", "18"}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			csv := strings.NewReader(tc.input)
			csvRows := ReadCSV(csv)

			if !reflect.DeepEqual(tc.want, csvRows) {
				t.Fatalf("expected %v, got: %v", tc.want, csvRows)
			}
		})
	}
}

// Happy path, because we will validate the CSV before passing it to this function.
func TestConvertRowToProblem_CorrectSliceFormat(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  Problem
	}{
		"simple": {
			input: []string{"1+1", "2"},
			want:  Problem{question: "1+1", answer: "2"},
		},
		"empty value for answer": {
			input: []string{"1+2", ""},
			want:  Problem{question: "1+2", answer: ""},
		},
		"empty value for question": {
			input: []string{"", "3"},
			want:  Problem{question: "", answer: "3"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := MapFields(tc.input)
			if !reflect.DeepEqual(tc.want, p) {
				t.Fatalf("expected: %v, got: %v", tc.want, p)
			}
		})
	}
}

func TestValidateCSVRows(t *testing.T) {
	tests := map[string]struct {
		input [][]string
		want  bool
	}{
		"simple": {input: [][]string{
			[]string{"1+1", "2"},
			[]string{"2+3", "5"},
		},
			want: true,
		},
		"invalid": {
			input: [][]string{
				[]string{"1+1", "2"},
				[]string{"2+3"},
			},
			want: false,
		}}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			valid := ValidateColumnCount(tc.input, 2)
			if tc.want != valid {
				t.Fatalf("expected: %v, actual: %v", tc.want, valid)
			}
		})
	}
}
