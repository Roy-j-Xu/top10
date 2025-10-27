import type { MessageManager } from "./message_manager";
import type { Message } from "./message_types";

export class SocketManager {
  private subscribers: Set<MessageManager> = new Set();
  private url: string = "";
  private ws!: WebSocket;
  private isConnected = false;

  private onMessage(event: MessageEvent) {
    const msg: Message = JSON.parse(event.data);
    this.subscribers.forEach(s => s.createHandler()(msg));
  }

  connect(url: string, ...msgManagers: MessageManager[]) {
    this.ws = new WebSocket(url);
    
    msgManagers.forEach(m => this.subscribe(m));
    this.ws.onmessage = this.onMessage;

    this.ws.onopen = () => {
      this.isConnected = true;
      console.log('Connected to server');
    };

    this.ws.onclose = () => {
      this.isConnected = false;
      console.log('Disconnected from server');
    };

    this.ws.onerror = (err) => {
      console.error('WebSocket error', err);
    };

    this.isConnected = true;
  }

  reconnect() {
    if (!this.isConnected || this.url === "") {
      return;
    }
    this.connect(this.url);
  }

  subscribe(msgManager: MessageManager) {
    this.subscribers.add(msgManager);
  }

  unsubscribe(msgManager: MessageManager) {
    this.subscribers.delete(msgManager);
  }

  send(msg: Message) {
    if (!this.isConnected) throw new Error('WebSocket not connected');
    this.ws.send(JSON.stringify(msg));
  }

  close() {
    this.ws.close();
    this.subscribers.forEach(s => this.unsubscribe(s));
    this.isConnected = false;
    this.url = "";
  }
}
