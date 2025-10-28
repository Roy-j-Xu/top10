import { useState, type ChangeEvent } from "react";
import gameService from "../core/game_service";
import type { RoomInfoResponse } from "../core";

export default function Main() {
  const [roomName, setRoomName] = useState("");
  const [playerName, setPlayerName] = useState("");
  const [roomInfo, setRoomInfo] = useState<RoomInfoResponse>();

  function handleInputChange(setter: (input: string) => void) {
    return (e: ChangeEvent<HTMLInputElement>) => {
      setter(e.target.value);
    };
  }

  const handleCreateRoom = async () => {
    try {
      await gameService.newGame(roomName, 4, "Top10");
      const info = await gameService.joinGame(roomName, playerName);
      setRoomInfo(info);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <input
        type="text"
        value={roomName}
        onChange={handleInputChange(setRoomName)}
        placeholder="room name"
      />
      <input
        type="text"
        value={playerName}
        onChange={handleInputChange(setPlayerName)}
        placeholder="player name"
      />
      <button onClick={handleCreateRoom}>create room</button>

      {roomInfo && (
        <div>
          <h3>Room info:</h3>
          <pre>{JSON.stringify(roomInfo, null, 2)}</pre>
        </div>
      )}
    </div>
  );
};