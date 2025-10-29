import { useEffect, useState } from "react"
import type { TopTenHandler, TurnInfoMsgData } from "../core/games/top10";
import type { RoomInfoResponse, SystemMessageHandler } from "../core";
import gameService from "../core/game_service";
import { useParams } from "react-router-dom";

export function TopTen() {
  const { roomName, playerName } = useParams();
  
  const gameMsgHandler = gameService.getHandler<TopTenHandler>("game");
  const systemMsgHandler = gameService.getHandler<SystemMessageHandler>("system");

  if (!roomName || !playerName) {
    return <h1>404 - Invalid room name or player name</h1>;
  }

  return (
    <GameBoard 
      roomName={roomName}
      playerName={playerName}
      gameMsgHandler={gameMsgHandler}
      sysMsgHandler={systemMsgHandler}
    />
  )
}

interface GameBoardParam {
  roomName: string;
  playerName: string;
  gameMsgHandler: TopTenHandler;
  sysMsgHandler: SystemMessageHandler;
}

function GameBoard(params: GameBoardParam) {
  const { roomName, playerName, gameMsgHandler, sysMsgHandler } = params;

  const [playerNum, setPlayerNum] = useState(1);
  const [roomInfo, setRoomInfo] = useState<RoomInfoResponse>();
  const [turnInfo, setTurnInfo] = useState<TurnInfoMsgData>();

  // update message handlers
  useEffect(() => {
    sysMsgHandler.onJoined((msg) => {
      console.log(msg.msg.playerName);
      setPlayerNum(prev => prev + 1);
    });
    sysMsgHandler.onLeft((msg) => {
      console.log(msg.msg.playerName);
      setPlayerNum(prev => prev - 1);
    })
    gameMsgHandler.onTurnInfo((msg) => {
      setTurnInfo(msg.msg);
    })
    window.addEventListener("beforeunload", gameService.endGame);
  }, [gameMsgHandler, sysMsgHandler]);
  
  // get room info
  useEffect(() => {
    gameService.getRoomInfo(roomName)
      .then(info => {
        setPlayerNum(info.players.length);
        setRoomInfo(info);
      });
  }, [roomName]);

  return (
    <div>
      <h1>Top10</h1>
      <div>Wellcome {playerName}</div>
      <div>{JSON.stringify(roomInfo)}</div>
      <div>Number of players: {playerNum}</div>
      <div>{JSON.stringify(turnInfo)}</div>
    </div>
  )
}