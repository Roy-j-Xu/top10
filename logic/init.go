package logic

var r *Room

func InitLogic() {
	r = NewRoom()
	LoadQuestionSet()
}
