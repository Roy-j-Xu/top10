import type { Game } from "./games/game"
import { MessageSender, SystemMessageHandler } from "./message"


export const registeredGames: Record<string, Game> = {
  "Top10": {
    url: "/ws",
    handlers: {"system": new SystemMessageHandler()},
    sender: MessageSender,
  },
}