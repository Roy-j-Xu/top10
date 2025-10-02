package core

var r *Room

func InitCore() {
	r = NewRoom()
	LoadQuestionSet()
}
