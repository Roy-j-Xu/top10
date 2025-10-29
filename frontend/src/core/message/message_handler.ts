import { SystemMsgType, type Message } from "./message_types";
import type { SocketManager } from "./socket_manager";
import { logMessage } from "./utils";

export class MessageHandler {
  protected handlers: Map<string, Set<(msg: Message) => void>> = new Map();

  createHandler(): (msg: Message) => void {
    return (msg: Message) => {
      const typeHandlers = this.handlers.get(msg.type);
      typeHandlers?.forEach(h => h(msg));
    };
  }

  register(msgType: string, handler: (msg: Message) => void): () => void {
    if (!this.handlers.has(msgType)) {
      this.handlers.set(msgType, new Set());
    }
    this.handlers.get(msgType)?.add(handler);
    return () => this.unregister(msgType, handler);
  }

  unregister(type: string, handler: (msg: Message) => void) {
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

  onJoined(handler: (msg: Message) => void) {
    this.register(SystemMsgType.JOINED, handler);
  }

  onLeft(handler: (msg: Message) => void) {
    this.register(SystemMsgType.LEFT, handler);
  }

  onStart(handler: (msg: Message) => void) {
    this.register(SystemMsgType.START, handler);
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

