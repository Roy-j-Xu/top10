import type { Game } from "./games/game"
import { MessageSender, SystemMessageHandler } from "./message"


export const registeredGames: Record<string, Game> = {
  "Top10": {
    minSize: 3,
    maxSize: 10,
    handlers: {"system": new SystemMessageHandler()},
    sender: MessageSender,
  },
}