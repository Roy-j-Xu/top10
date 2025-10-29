import type { MessageHandler } from "./message_handler";
import type { Message } from "./message_types";

export class SocketManager {
  private subscribers: Set<MessageHandler> = new Set();
  private ws!: WebSocket;

  public isConnected = false;

  private onMessage(event: MessageEvent) {
    const msg: Message = JSON.parse(event.data);
    this.subscribers.forEach(s => s.createHandler()(msg));
  }

  connect(url: string, onClose: () => void, ...msgManagers: MessageHandler[]) {
    if (this.isConnected) {
      return;
    }
    if (!url) {
      throw Error("socket reconnect: url not known");
    }

    this.ws = new WebSocket(url);
    
    msgManagers.forEach(m => this.subscribe(m));
    this.ws.onmessage = this.onMessage.bind(this);

    this.ws.onopen = () => {
      this.isConnected = true;
      console.log('connected to server');
    };

    this.ws.onclose = onClose;

    this.ws.onerror = (err) => {
      console.error('webSocket error', err);
    };

    this.isConnected = true;
  }

  subscribe(msgManager: MessageHandler) {
    this.subscribers.add(msgManager);
  }

  unsubscribe(msgManager: MessageHandler) {
    this.subscribers.delete(msgManager);
  }

  send(msg: Message | string) {
    if (!this.isConnected) throw new Error('WebSocket not connected');
    this.ws.send(JSON.stringify(msg));
  }

  close() {
    if (!this.isConnected) {
      return;
    }
    this.ws.close();
    this.subscribers.forEach(s => this.unsubscribe(s));
    this.isConnected = false;
  }
}
