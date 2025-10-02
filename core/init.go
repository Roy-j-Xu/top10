package core

var r *Room

func InitCore(msgrs []Messager) {
	r = NewRoom(msgrs)
	LoadQuestionSet()
}
