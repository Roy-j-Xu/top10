package logic

import (
	"bufio"
	"math/rand/v2"
	"os"
)

var QuestionSet []string

func LoadQuestionSet() {
	file, err := os.Open("./question_set.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if text := scanner.Text(); text != "" {
			QuestionSet = append(QuestionSet, text)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func RandomQuestion() string {
	return QuestionSet[rand.IntN(len(QuestionSet))]
}
