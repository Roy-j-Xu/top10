package tests

import (
	"testing"
	"top10/logic"
)

func TestRandomQuestion(t *testing.T) {
	logic.LoadQuestionSet()
	print(logic.RandomQuestion())
}
