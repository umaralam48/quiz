package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	pfilename := flag.String("csv", "problems.csv", " The CSV file to read the problems from in the \"problem, answer\" format.")
	limit := flag.Int("limit", 30, "Time limit for the quiz.")
	flag.Parse()
	file, err := os.Open(*pfilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s !\n", *pfilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the CSV file : %s", *pfilename))
	}
	problems := parseCSV(lines)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	correct := 0
	ansch := make(chan string)
	for i, problem := range problems {
		fmt.Printf("Question %d : %s =\n", i+1, problem.q)
		go func() {
			var ans string
			fmt.Scanf("%s", &ans)
			ansch <- ans
		}()
		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d\n", correct, len(problems))
			return
		case ans := <-ansch:
			if ans == problem.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func parseCSV(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, p := range lines {
		problems[i] = problem{
			q: p[0],
			a: p[1],
		}
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
