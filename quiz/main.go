package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getInput(reader *bufio.Scanner) (string, error) {
	var str string
	_, err := fmt.Scanln(&str)
	if err != nil {
		return "", err
	}
	str = strings.TrimRight(str, "\n")
	return str, nil
}

func readCSV(filename string) ([][]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(strings.NewReader(string(bytes)))

	var allQuestions [][]string

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		allQuestions = append(allQuestions, record)
	}

	return allQuestions, nil
}

func main() {
	reader := bufio.NewScanner(os.Stdin)

	var filename string
	var limit int
	flag.StringVar(&filename, "file", "problem.csv", "File with CSV format to read the problems from")
	flag.IntVar(&limit, "timer", 30, "Countdown timer")

	flag.Parse()

	questions, err := readCSV(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Enter when you're ready!")
	getInput(reader)

	fmt.Println("Timer starts!")

	timer1 := time.After(time.Duration(limit) * time.Second)
	correct := 0

loop:
	for index, question := range questions {
		fmt.Println("Question " + strconv.Itoa(index+1))
		fmt.Println(question[0])
		answerChannel := make(chan string, 1)
		go func() {
			input, err := getInput(reader)
			if err != nil {
				log.Fatalln(err.Error())
			}
			answerChannel <- input
		}()

		select {
		case <-timer1:
			break loop
		case answer := <-answerChannel:
			if answer == question[1] {
				fmt.Println("Answer is correct!")
				correct++
			} else {
				fmt.Println("Answer wrong!")
				fmt.Println("You answered: " + answer)
				fmt.Println("Correct: " + question[1])
				break
			}
			break
		}
	}

	fmt.Printf("You correctly answered %d out of %d questions", correct, len(questions))
}
