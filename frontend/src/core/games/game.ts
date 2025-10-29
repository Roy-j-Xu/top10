import type { MessageHandler, MessageSender, SocketManager } from "../message";

export interface Game {
  minSize: number;
  maxSize: number;
  handlerFactories: Record<string, () => MessageHandler>;
  senderFactory: (s: SocketManager) => MessageSender;
}