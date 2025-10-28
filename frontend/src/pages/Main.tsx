import { useState, type ChangeEvent } from "react";
import gameService from "../core/game_service";
import type { RoomInfoResponse } from "../core";

export default function Main() {
  const [inputValue, setInputValue] = useState("");
  const [roomInfo, setRoomInfo] = useState<RoomInfoResponse>();

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleButtonClick = async () => {
    try {
      const resp = await gameService.newGame(inputValue, 4, "Top10");
      setRoomInfo(resp);
    } catch (error) {
      console.error("error:", error);
    }
  };

  return (
    <div>
      <input
        type="text"
        value={inputValue}
        onChange={handleInputChange}
        placeholder="room name"
      />
      <button onClick={handleButtonClick}>create room</button>

      {roomInfo && (
        <div>
          <h3>Room info:</h3>
          <pre>{JSON.stringify(roomInfo, null, 2)}</pre>
        </div>
      )}
    </div>
  );
};