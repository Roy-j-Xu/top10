import { useEffect, useState } from "react"
import type { TopTenHandler, TurnInfoResponse } from "../core/games/top10";
import type { RoomInfoResponse, SystemMessageHandler } from "../core";
import gameService from "../core/game_service";
import { useParams } from "react-router-dom";

export function TopTen() {
  const { roomName } = useParams();
  
  const [playerNum, setPlayerNum] = useState(1);
  const [roomInfo, setRoomInfo] = useState<RoomInfoResponse>();
  const [turnInfo, setTurnInfo] = useState<TurnInfoResponse>();

  const gameMsgHandler = gameService.getHandler<TopTenHandler>("game");
  const systemMsgHandler = gameService.getHandler<SystemMessageHandler>("system");

  // update message handlers
  useEffect(() => {
    systemMsgHandler.onJoined(() => {
      setPlayerNum(prev => prev + 1);
    });
    systemMsgHandler.onLeft(() => {
      setPlayerNum(prev => prev - 1);
    })
    gameMsgHandler.onTurnInfo((msg) => {
      setTurnInfo(msg.msg as TurnInfoResponse);
    })
    window.addEventListener("beforeunload", gameService.endGame);
  }, [gameMsgHandler, systemMsgHandler]);
  
  // update room info
  useEffect(() => {
    if (!roomName) {
      return;
    }
    gameService.getRoomInfo(roomName)
      .then(info => {
        setPlayerNum(info.players.length);
        setRoomInfo(info);
      });
  }, [playerNum, roomName]);

  if (!roomName) {
    return <h1>404 - Invalid room name</h1>;
  }
  if (!gameMsgHandler || ! systemMsgHandler) {
    return <h1>500 - Internal server error</h1>
  }

  return (
    <div>
      <div>{JSON.stringify(roomInfo)}</div>
      <div>{playerNum}</div>
      <div>{JSON.stringify(turnInfo)}</div>
    </div>
  )
}