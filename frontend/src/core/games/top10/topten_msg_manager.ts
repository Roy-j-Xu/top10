import { logMessage, MessageHandler, MessageSender, type HandlerFunc } from "../../message";
import { TopTenMsgType, TopTenPlayerMsgType, type GameInfo } from "./types";

export class TopTenHandler extends MessageHandler {
  constructor() {
    super();
    this.useLogger();
  }

  onGameInfo(handler: HandlerFunc<GameInfo>) {
    this.register(TopTenMsgType.GAME_INFO, handler as HandlerFunc);
  }

  onStartGuessing(handler: HandlerFunc<GameInfo>) {
    this.register(TopTenMsgType.START_GUESSING, handler as HandlerFunc);
  }

  onStart(handler: HandlerFunc<GameInfo>) {
    this.register(TopTenMsgType.START, handler as HandlerFunc);
  }

  onFinished(handler: HandlerFunc) {
    this.register(TopTenMsgType.FINISHED, handler);
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