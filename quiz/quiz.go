package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Problem represents a question/answer pair
type Problem struct {
	Question string
	Answer   string
}

// LoadProblems unmarshals a two column CSV file into a slice of Problems
func LoadProblems(path string) []Problem {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Couldn't open file %s: %s", path, err)
	}
	defer f.Close()

	csv := csv.NewReader(f)
	lines, err := csv.ReadAll()
	if err != nil {
		log.Fatalf("Couldn't read file %s: %s", path, err)
	}

	problems := make([]Problem, len(lines))
	for i, line := range lines {
		problems[i] = Problem{
			Question: line[0],
			Answer:   line[1],
		}
	}

	return problems
}

// ShuffleProblems pseudo-randomizes the order of problems.
func ShuffleProblems(problems []Problem) {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

// Quiz asks the user interactively to answer questions specified in problems.
// It returns the number of correctly answered questions once the user has
// answered all of the questions or timeLimit seconds has elapsed.
// Ignores capitalization and leading/trailing whitespace when comparing answers.
func Quiz(problems []Problem, timeLimit int) (correct int) {
	ch := make(chan *string)
	quit := make(chan int)

	go func() {
		<-time.After(time.Duration(timeLimit) * time.Second)
		quit <- 1
	}()

	for i, problem := range problems {
		answer := new(string)
		go func() {
			input := new(string)
			fmt.Printf("Problem #%d: %s? ", i+1, problem.Question)
			fmt.Scanf("%s", input)
			ch <- input
		}()
		select {
		case answer = <-ch:
			if formatAnswer(*answer) == formatAnswer(problem.Answer) {
				correct++
			}
		case <-quit:
			return
		}
	}
	return
}

func formatAnswer(answer string) string {
	return strings.TrimSpace(strings.ToLower(answer))
}

func main() {
	problemFile := flag.String("i", "problems.csv", "Input problem CSV file")
	timeLimit := flag.Int("t", 30, "Time limit in seconds")
	shuffle := flag.Bool("R", false, "Randomize the order of the problems")
	flag.Parse()

	problems := LoadProblems(*problemFile)

	rand.Seed(time.Now().UnixNano())
	if *shuffle == true {
		ShuffleProblems(problems)
	}

	fmt.Printf("You'll have %d seconds to answer %d questions\n", *timeLimit, len(problems))
	fmt.Println("Press enter to begin")
	fmt.Scanln()

	correct := Quiz(problems, *timeLimit)
	fmt.Printf("\nYou got %d/%d correct\n", correct, len(problems))
}
