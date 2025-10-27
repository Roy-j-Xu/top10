import { SystemMsgType } from "./message_types";
import { MessageManager } from "./message_manager";


export class SystemMessageManager extends MessageManager {
  constructor(useLogger=false) {
    super();
    if (useLogger) {
      this.useLogger()
    }
  }

  private useLogger() {
    for (const type in SystemMsgType) {
      if (type === SystemMsgType.S_ERROR) {
        this.register(type, (msg) => console.error(msg));
      }
      this.register(type, (msg) => console.log(msg));
    }
  }

}
