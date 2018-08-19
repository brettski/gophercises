package main

import (
	"flag"
	"os"
	"fmt"
	"encoding/csv"
	"strings"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv in the format of 'question, answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("failed to open the CSV file: %s\n", *csvFilename)
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	
	problems := parseLines(lines)
	
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i + 1, p.q)
		var answer string
		// fmt.Scanf good enough for this as it is single value entries
		fmt.Scanf("%s\n", &answer)
		if answer == p.a {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	// not useing append as we know the size of data
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}