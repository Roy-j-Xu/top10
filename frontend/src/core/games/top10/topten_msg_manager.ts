import { logMessage, MessageHandler, MessageSender, type Message } from "../../message";
import { TopTenMsgType, TopTenPlayerMsgType } from "./types";

export class TopTenHandler extends MessageHandler {
  constructor() {
    super();
    this.useLogger();
  }

  onTurnInfo(handler: (msg: Message) => void) {
    this.register(TopTenMsgType.TURN_INFO, handler);
  }

  onAssignNumbers(handler: (msg: Message) => void) {
    this.register(TopTenMsgType.ASSIGN_NUMBERS, handler);
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
  ready() {
    this.send({
      type: TopTenPlayerMsgType.READY,
      msg: "ready",
    });
  }

  setQuestion(question: string) {
    this.send({
      type: TopTenPlayerMsgType.SET_QUESTION,
      msg: question,
    });
  }

  chooseOrder(num: number) {
    this.send({
      type: TopTenPlayerMsgType.CHOOSE_ORDER,
      msg: num,
    });
  }
}