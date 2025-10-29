package room

type Message struct {
	Type string `json:"type"`
	Msg  any    `json:"msg"`
}

type SystemMsgType string

const (
	// Messages from system to players
	S_JOINED    SystemMsgType = "system:joined"
	S_LEFT      SystemMsgType = "system:left"
	S_READY     SystemMsgType = "systen:ready"
	S_START     SystemMsgType = "system:start"
	S_BROADCAST SystemMsgType = "system:broadcast"
	S_ERROR     SystemMsgType = "system:error"

	// Messages from players to system
	SP_READY SystemMsgType = "system-player:ready"
	SP_LEFT  SystemMsgType = "system-player:leave"
)

type playerMsgData struct {
	PlayerName string `json:"playerName"`
	Message    string `json:"message"`
}

type JoinedMsgData playerMsgData
type LeftMsgData playerMsgData
type ReadyMsgData playerMsgData

func SystemMsgOf(msgType SystemMsgType, msg any) Message {
	return Message{
		Type: string(msgType),
		Msg:  msg,
	}
}

func JoinedMsgOf(playerName string, msg string) Message {
	return Message{
		Type: string(S_JOINED),
		Msg: JoinedMsgData{
			PlayerName: playerName,
			Message:    msg,
		},
	}
}

func LeftMsgOf(playerName string, msg string) Message {
	return Message{
		Type: string(S_LEFT),
		Msg: LeftMsgData{
			PlayerName: playerName,
			Message:    msg,
		},
	}
}

func ReadyMsgOf(playerName string, msg string) Message {
	return Message{
		Type: string(S_LEFT),
		Msg: ReadyMsgData{
			PlayerName: playerName,
			Message:    msg,
		},
	}
}
