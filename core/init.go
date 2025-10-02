package core

var r *Room

func InitCore(msgr Messager) {
	r = NewRoom(msgr)
	LoadQuestionSet()
}
