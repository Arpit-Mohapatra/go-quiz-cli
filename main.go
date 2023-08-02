package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the form of questions,answers")
	timeLimit := flag.Int("limit", 30, "time limit to complete the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("Could not open %v", *csvFile))
	}
	defer file.Close()
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Could not read file")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem %d is %s = ", i +1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <-answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou got %d correct out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You got %d correct out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return problems
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}