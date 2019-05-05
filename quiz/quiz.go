package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filePtr := flag.String("file", "problems.csv", "file name of the question set")
	timeLimitPtr := flag.Int("timelimit", 30, "Time limit for the quiz in seconds, defaults to 30 i.e. no time limit")
	flag.Parse()
	f, err := os.Open(*filePtr)
	checkError(err)
	defer f.Close()

	csv := csv.NewReader(f)
	questions, err := csv.ReadAll()
	checkError(err)
	quiz := Quiz{questionCount: len(questions), questions: questions}
	select {
	case <-quiz.startQuiz():
		quiz.endQuiz()
	case <-time.After(time.Duration(*timeLimitPtr) * time.Second):
		quiz.timeOut()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Quiz represents a quiz read from a csv file
type Quiz struct {
	score         int
	questionCount int
	questions     [][]string
}

func (q *Quiz) startQuiz() chan bool {
	ch := make(chan bool)
	go func() {
		for {
			for _, row := range q.questions {
				question, correctAnswer := string(row[0]), string(row[1])
				fmt.Printf("What is %s?\n", question)
				var userInput string
				fmt.Scanln(&userInput)
				userInput = strings.ToLower(strings.TrimSpace(userInput))
				if userInput == correctAnswer {
					q.score++
				}
			}
		}
		ch <- true
	}()
	return ch
}

func (q *Quiz) endQuiz() {
	fmt.Printf("Quiz ended! You scored %d/%d!", q.score, q.questionCount)
}

func (q *Quiz) timeOut() {
	fmt.Printf("Time's up! You scored %d/%d!", q.score, q.questionCount)
}