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

type RoomInfo struct {
	RoomName string   `json:"roomName"`
	RoomSize int      `json:"roomSize"`
	Game     string   `json:"game"`
	Players  []string `json:"players"`
	InGame   bool     `json:"inGame"`
}

type playerMsgData struct {
	PlayerName string   `json:"playerName"`
	RoomInfo   RoomInfo `json:"roomInfo"`
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

func JoinedMsgOf(playerName string, roomInfo RoomInfo) Message {
	return Message{
		Type: string(S_JOINED),
		Msg: JoinedMsgData{
			PlayerName: playerName,
			RoomInfo:   roomInfo,
		},
	}
}

func LeftMsgOf(playerName string, roomInfo RoomInfo) Message {
	return Message{
		Type: string(S_LEFT),
		Msg: LeftMsgData{
			PlayerName: playerName,
			RoomInfo:   roomInfo,
		},
	}
}

func ReadyMsgOf(playerName string, roomInfo RoomInfo) Message {
	return Message{
		Type: string(S_READY),
		Msg: LeftMsgData{
			PlayerName: playerName,
			RoomInfo:   roomInfo,
		},
	}
}
