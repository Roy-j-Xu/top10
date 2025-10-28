import type { MessageHandler, MessageSender, SocketManager } from "../message";

export interface Game {
  url: string;
  handlers: Record<string, MessageHandler>;
  sender: new (socket: SocketManager) => MessageSender;
}