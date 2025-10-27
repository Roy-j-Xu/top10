export interface Message {
  type: string;
  msg: unknown;
}

export function messageNamespace(msg: Message) {
	return msg.type.split(":")[0];
}

export type MessageHandler = (msg: Message) => void;

export const SystemMsgType = {
  S_JOINED: "system:joined",
	S_LEFT: "system:left",
	S_START: "system:start",
	S_BROADCAST: "system:broadcast",
	S_ERROR: "system:error",

  SP_READY: "system-player:ready",
	SP_LEFT: "system-player:leave",
} as const;

export type SystemMsgType = typeof SystemMsgType[keyof typeof SystemMsgType];
