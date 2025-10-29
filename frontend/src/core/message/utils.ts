import type { Message } from "./message_types";

export function logMessage(msg: Message) {
  console.log(`[${msg.type}] ${msg.msg}`)
}