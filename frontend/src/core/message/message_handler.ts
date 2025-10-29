import { SystemMsgType, type JoinedMsgData, type LeftMsgData, type Message, type ReadyMsgData } from "./message_types";
import type { SocketManager } from "./socket_manager";
import { logMessage } from "./utils";

export type HandlerFunc<T=unknown> = (msg: Message<T>) => void;

export class MessageHandler {
  protected handlers: Map<string, Set<HandlerFunc>> = new Map();

  createHandler(): (msg: Message) => void {
    return (msg: Message) => {
      const typeHandlers = this.handlers.get(msg.type);
      typeHandlers?.forEach(h => h(msg));
    };
  }

  register(msgType: string, handler: HandlerFunc): () => void {
    if (!this.handlers.has(msgType)) {
      this.handlers.set(msgType, new Set());
    }
    this.handlers.get(msgType)?.add(handler);
    return () => this.unregister(msgType, handler);
  }

  unregister(type: string, handler: HandlerFunc) {
    const handlerSet = this.handlers.get(type);
    if (!handlerSet) return;
    handlerSet.delete(handler);
    if (handlerSet.size === 0) this.handlers.delete(type);
  }

}

export class MessageSender {
  private socket: SocketManager;

  constructor(socket: SocketManager) {
    this.socket = socket;
  }
  
  send<T extends Message>(msg: T) {
    this.socket.send(msg);
  }
}

export class SystemMessageHandler extends MessageHandler {
  constructor(useLogger=false) {
    super();
    if (useLogger) {
      this.useLogger()
    }
  }

  onJoined(handler: (msg: Message<JoinedMsgData>) => void) {
    this.register(SystemMsgType.JOINED, handler as HandlerFunc);
  }

  onLeft(handler: (msg: Message<LeftMsgData>) => void) {
    this.register(SystemMsgType.LEFT, handler as HandlerFunc);
  }

  onReady(handler: (msg: Message<ReadyMsgData>) => void) {
    this.register(SystemMsgType.LEFT, handler as HandlerFunc);
  }

  onStart(handler: (msg: Message) => void) {
    this.register(SystemMsgType.START, handler as HandlerFunc);
  }

  private useLogger() {
    Object.values(SystemMsgType).forEach(type => {
      if (type === SystemMsgType.ERROR) {
        this.register(type, (msg) => console.error(msg));
      }
      this.register(type, (msg) => logMessage(msg));
    });
  }
}

