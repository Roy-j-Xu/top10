package game

import (
	"bufio"
	"math/rand/v2"
	"os"
)

var QuestionSet []string

func LoadQuestionSet() {
	file, err := os.Open("resource/question_set.txt")
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

func RandomQuestions(num int) []string {
	var indices []int

	for len(indices) < num {
		n := rand.IntN(len(QuestionSet))

		noRepeat := true
		for _, number := range indices {
			if n == number {
				noRepeat = false
				break
			}
		}

		if noRepeat {
			indices = append(indices, n)
		}
	}

	result := make([]string, num)
	for i, index := range indices {
		result[i] = QuestionSet[index]
	}

	return result
}
