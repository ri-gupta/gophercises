package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func readCSVFile(filename *string) ([]byte, error) {
	f, err := os.Open(*filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f) // load bytes data in memory

	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseCsvData(data []byte) *csv.Reader {
	reader := csv.NewReader(bytes.NewReader(data))

	return reader
}

type Quiz struct {
	questions [][]string
}

func (q *Quiz) start(score *int, timer *time.Timer) {
	answerchan := make(chan string)

	for i, record := range q.questions {

		fmt.Print("# Question ", i+1, " : ", record[0], " ")
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answer = strings.TrimSpace(answer)

			answerchan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			return
		case answer := <-answerchan:
			if answer == record[1] {
				*score++
			}
		}
	}
}

func main() {

	// cmd parameter filename with default value problems.csv
	filename := flag.String("filename", "problems.csv", "File used to serve questions for the quiz!")
	timer := flag.Int("timer", 30, "Number of seconds given to complete the quiz!")

	flag.Parse()

	rawData, err := readCSVFile(filename)

	if err != nil {
		fmt.Println("Error reading csv: ", err)
		os.Exit(1)
	}

	parsedData := parseCsvData(rawData)

	questions, err := parsedData.ReadAll()
	totalQuestions := len(questions)
	score := 0

	quiz := Quiz{questions: questions}

	if err != nil {
		fmt.Println("Error reading questions data: ", err)
		os.Exit(1)
	}

	fmt.Println("Press Enter to start the quiz: ")
	stdReader := bufio.NewReader(os.Stdin)
	_, err = stdReader.ReadString('\n')

	if err != nil {
		fmt.Println("Invalid input, Please enter any key & press enter")
	}

	// starts the quiz timer
	quizTimer := time.NewTimer(time.Duration(*timer) * time.Second)
	quiz.start(&score, quizTimer)

	fmt.Println("You scored", score, "out of", totalQuestions)
}
