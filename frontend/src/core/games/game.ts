import type { MessageHandler, MessageSender, SocketManager } from "../message";

export interface Game {
  minSize: number;
  maxSize: number;
  handlers: Record<string, MessageHandler>;
  sender: new (socket: SocketManager) => MessageSender;
}