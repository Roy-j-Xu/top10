import type { Game } from "./games/game"
import { TopTenHandler, TopTenSender } from "./games/top10"
import { SystemMessageHandler } from "./message"


export const registeredGames: Record<string, Game> = {
  "Top10": {
    minSize: 3,
    maxSize: 10,
    handlerFactories: {
      "system": () => new SystemMessageHandler(true),
      "game": () => new TopTenHandler(),
    },
    senderFactory: (s) => new TopTenSender(s),
  },
}