import type { Game } from "./games/game";
import { MessageHandler, MessageSender, SocketManager, type RoomInfoResponse } from "./message";
import { registeredGames } from "./registered_games";
import config from "../../config.json"

class GameService {
  private socketManager = new SocketManager(); 
  private sender?: MessageSender;
  private game?: Game;
  
  async newGame(roomName: string, roomSize: number, gameName: string): Promise<RoomInfoResponse> {
    const game = registeredGames[gameName];
    if (roomSize < game.minSize || game.maxSize < roomSize) {
      throw Error("invalid room size for game");
    }

    const resp = await this.newRoom(roomName, roomSize);

    this.sender = new game.sender(this.socketManager);
    this.game = game
    return resp;
  }

  private async newRoom(roomName: string, roomSize: number): Promise<RoomInfoResponse> {
    const body = JSON.stringify({ roomName: roomName, roomSize: roomSize});
    console.log(body);
    const response = await fetch(`${config["apiUrl"]}/create-room`, {
      method: "POST",
      body: body,
    });
    if (!response.ok) {
      throw Error(`creating room: ${JSON.stringify(response)}`);
    }
    return await response.json();
  }

  async joinGame(roomName: string, playerName: string): Promise<RoomInfoResponse> {
    const data = await this.getRoomInfo(roomName);
    this.game = registeredGames[data.game]
    this.socketManager.connect(
      `${config["socketUrl"]}?roomName=${roomName}&playerName=${playerName}`,
      ...Object.values(this.getHandlers()),
    );
    return data;
  }

  async getRoomInfo(roomName: string): Promise<RoomInfoResponse> {
    const response = await fetch(`${config["apiUrl"]}/room-info?roomName=${roomName}`);
    if (!response.ok) {
      throw Error(`creating room: ${JSON.stringify(response)}`);
    }
    return await response.json();
  }
  
  endGame() {
    this.game = undefined;
    this.sender = undefined;
    this.socketManager.close();
  }

  getSender<S extends MessageSender>(): S {
    if (this.game === undefined) {
      throw Error("game not found, unable to get sender");
    }
    return this.sender as S;
  }

  getHandlers(): Record<string, MessageHandler> {
    if (this.game === undefined) {
      throw Error("game not found, unable to get handler");
    }
    return this.game.handlers;
  }

}

const gameService = new GameService();

export default gameService;