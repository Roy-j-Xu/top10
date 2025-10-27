import type { Message, MessageHandler } from "./message_types";

export class MessageManager {
  protected handlers: Map<string, Set<MessageHandler>> = new Map();

  createHandler(): MessageHandler {
    return (msg: Message) => {
      const typeHandlers = this.handlers.get(msg.type);
      typeHandlers?.forEach(h => h(msg));
    };
  }

  register(msgType: string, handler: MessageHandler): () => void {
    if (!this.handlers.has(msgType)) {
      this.handlers.set(msgType, new Set());
    }
    this.handlers.get(msgType)?.add(handler as MessageHandler);
    return () => this.unregister(msgType, handler as MessageHandler);
  }

  unregister(type: string, handler: MessageHandler) {
    const handlerSet = this.handlers.get(type);
    if (!handlerSet) return;
    handlerSet.delete(handler);
    if (handlerSet.size === 0) this.handlers.delete(type);
  }

}