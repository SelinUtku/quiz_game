package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type quiz struct {
	q string
	a string
}

func readFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	all, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return all
}

func initializeStruct(all [][]string) []quiz {
	setOfQuiz := make([]quiz, len(all))
	for i, line := range all {
		setOfQuiz[i] = quiz{
			q: line[0],
			a: strings.Trim(line[1], " "),
		}

	}
	return setOfQuiz
}

func startGame(setOfQuiz []quiz, limit int) int {
	var correct int = 0
	answerC := make(chan string)
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	for _, line := range setOfQuiz {
		go func(line quiz) {
			var input string
			fmt.Print(line.q, " = ")
			fmt.Scanf("%s", &input)
			answerC <- input
		}(line)
		select {
		case <-timer.C:
			fmt.Println("\nTime is up !", correct, "out of", len(setOfQuiz))
			os.Exit(0)
		case answer := <-answerC:
			timer.Reset(time.Duration(limit) * time.Second)
			if answer == line.a {
				correct++
			}
		}
	}
	return correct
}

func main() {

	fileName := flag.String("file", "problems.csv", "Provide a csv file Format:[Question,Answer]")
	limit := flag.Int("limit", 5, "Time limit in seconds")
	flag.Parse()
	setOfQuiz := initializeStruct(readFile(*fileName))

	fmt.Println("Quiz game is working on", *fileName)
	correct := startGame(setOfQuiz, *limit)
	fmt.Println(correct, "out of", len(setOfQuiz))
}
