package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	duration := flag.Int("time", 25, "duration for timeout")
	fileName := flag.String("file", "problems.csv", "The CSV file to use for the quiz")
	records := readFileForQuiz(*fileName)

	flag.Parse()
	timeout := time.After(time.Duration(*duration) * time.Second)
	stop := make(chan struct{})
	done := make(chan struct{})
	count := 0

	go func() {
		quiz(records, stop, done, &count)
	}()

	select {
	case <-timeout:
		fmt.Println("Timeout reached, stopping execution")
		close(stop)
	case <-done:
		fmt.Println("Quiz completed")
	}

	fmt.Println("Out of ", len(records), "questions, You entered ", count, " correct answers")

}

func readFileForQuiz(fileName string) [][]string {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error while reading csv %s", err)
	}
	return records
}

func quiz(records [][]string, stop chan struct{}, done chan struct{}, count *int) {

	defer close(done)
	for _, record := range records {
		select {
		case <-stop:
			fmt.Println("Quiz stopped due to timeout")
			return

		default:
			fmt.Println(record[0])
			input := readInputFromUser()
			if input == record[1] {
				*count++
			}
		}
	}
}

func readInputFromUser() string {
	readerFromUser := bufio.NewReader(os.Stdin)
	input, _ := readerFromUser.ReadString('\n')
	return strings.TrimSpace(input)
}
