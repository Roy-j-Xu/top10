import { logMessage, MessageHandler, MessageSender } from "../../message";
import { TopTenMsgType } from "./message_types";

export class TopTenHandler extends MessageHandler {
  constructor() {
    super();
    this.useLogger();
  }

  useLogger() {
    Object.values(TopTenMsgType).forEach(type => {
      if (type === TopTenMsgType.ERROR) {
        this.register(type, (msg) => console.error(msg));
      }
      this.register(type, (msg) => logMessage(msg));
    });
  }
}

export class TopTenSender extends MessageSender {

}