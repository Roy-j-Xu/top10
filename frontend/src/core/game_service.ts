import type { Game } from "./games/game";
import { MessageHandler, MessageSender, SocketManager } from "./message";
import { registeredGames } from "./registered_games";
import config from "../../config.json"

class GameService {
  private socketManager = new SocketManager(); 
  private sender?: MessageSender;
  private game?: Game;
  
  runGame(gameName: string) {
    this.game = registeredGames[gameName];
    this.sender = new this.game.sender(this.socketManager);

    const url = `${config['baseUrl']}${this.game.url}}`;
    this.socketManager.connect(url, ...Object.values(this.game.handlers));
  }
  
  endCurrentGame() {
    this.game = undefined;
    this.sender = undefined;
    this.socketManager.close();
  }

  getSender<S extends MessageSender>(): S {
    if (this.game === undefined) {
      throw Error("game not in play, unable to get sender");
    }
    return this.sender as S;
  }

  getHandlers(): Record<string, MessageHandler> {
    if (this.game === undefined) {
      throw Error("game not in play, unable to get handler");
    }
    return this.game.handlers;
  }

}

const gameService = new GameService();

export default gameService;