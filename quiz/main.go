package main

/* Give the user a quiz over the CLI, where the questions and answers
 * are defined in a CSV file
 */

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// Retrieving command line arguments
	csvFileName := flag.String("csv", "quiz.csv", "A CSV file with the format: question, answer")
	limit := flag.Int("limit", 3, "Time limit for each question in seconds")
	flag.Parse()

	// Retrieve the problem set from the specified file
	problems := make(problemSet, 0)
	problems.PopulateFromCSV(*csvFileName)

	timeLimit := time.Duration(*limit) * time.Second

	println("The quiz is starting!")

	wrongAns := make([]int, 0, cap(problems))
	missedAns := make([]int, 0, cap(problems))
	correct := 0

	tmr := time.NewTimer(timeLimit)
	for i, p := range problems { // Give the quiz to the user

		// Set up the question an prompt the user for an answer
		fmt.Printf("Problem %d\t%s = ", i+1, p.ques)
		ansCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ansCh <- strings.TrimSpace(answer)
		}()

		// Reset the timer
		tmr.Stop()
		tmr.Reset(timeLimit)

		// Wait until the answer is received or the timer expires
		select {
		case ans := <-ansCh:
			if ans == p.ans {
				correct++
			} else {
				wrongAns = append(wrongAns, i)
			}
		case <-tmr.C:
			fmt.Printf("Out of time!\n")
			missedAns = append(missedAns, i)
		}
	}

	fmt.Printf("\nYou got %d out of %d questions correct.\n", correct, len(problems))

	// Show information on the questions the user got wrong or missed
	if len(wrongAns) > 0 {
		println("\nSolutions to the ones you got wrong:")
		for _, val := range wrongAns {
			fmt.Printf("Problem %d: %s = %s\n", val+1, problems[val].ques, problems[val].ans)
		}
	}
	if len(missedAns) > 0 {
		println("\nSolutions of the ones you didn't answer on time:")
		for _, val := range missedAns {
			fmt.Printf("Problem %d: %s = %s\n", val+1, problems[val].ques, problems[val].ans)
		}
	}
}

type problem struct {
	ques string
	ans  string
}

type problemSet []problem

func (problems *problemSet) PopulateFromCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil { // Could not open the file
		println("Could not open file", fileName)
		os.Exit(1)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil { // The CSV could not be parsed
		println("The CSV file could not be parsed")
		os.Exit(1)
	}

	for _, line := range lines { // Convert the lines into questions and answers
		*problems = append(*problems, problem{
			ques: strings.TrimSpace(line[0]),
			ans:  strings.TrimSpace(line[1]),
		})
	}

	if file.Close() != nil { // Could not close the file
		println("File", fileName, "could not be closed")
		os.Exit(1)
	}
}
