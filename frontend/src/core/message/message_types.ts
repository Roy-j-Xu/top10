export interface Message<T = unknown> {
  type: string;
  msg: T;
}

export function messageNamespace(msg: Message) {
	return msg.type.split(":")[0];
}

export const SystemMsgType = {
  JOINED: "system:joined",
  LEFT: "system:left",
  READY: "system:ready",
  START: "system:start",
  BROADCAST: "system:broadcast",
  ERROR: "system:error",
} as const;

export const SystemPlayerMsgType = {
  READY: "system-player:ready",
  LEFT: "system-player:leave",
} as const;

export interface RoomInfoResponse {
  roomName: string,
	roomSize: number,
	game:     string,
	players:  string[],
  inGame:   boolean,
}

interface PlayerMsgData {
  playerName: string;
  message: string;
}

export type {
  PlayerMsgData as JoinedMsgData,
  PlayerMsgData as LeftMsgData,
  PlayerMsgData as ReadyMsgData,
};