package main

import "top10/logic"

func main() {
	logic.InitLogic()
	print(logic.RandomQuestion())
}
