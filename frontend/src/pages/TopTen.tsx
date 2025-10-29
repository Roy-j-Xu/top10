import { useEffect, useState } from "react"
import type { TopTenHandler, TurnInfoMsgData } from "../core/games/top10";
import type { RoomInfoResponse, SystemMessageHandler } from "../core";
import gameService from "../core/game_service";
import { useParams } from "react-router-dom";

export function TopTen() {
  const { roomName, playerName } = useParams();

  const [connected, setConnected] = useState<boolean>(false);

  // reconnect for mobil platforms
  useEffect(() => {
    if (!roomName || !playerName) { return; }

    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        if (!gameService.isJoined()) {
          gameService.joinGame(roomName, playerName);
        }
      }
    };

    document.addEventListener("visibilitychange", handleVisibilityChange);
    return () => document.removeEventListener("visibilitychange", handleVisibilityChange);
  }, [playerName, roomName]);

  // join game if not already in
  useEffect(() => {
    if (!roomName || !playerName) { return; }
    if (!gameService.isJoined()) {
      gameService.joinGame(roomName, playerName)
        .then(() => {
          setConnected(true);
          window.addEventListener("beforeunload", gameService.leaveGame);
        })
    } else {
      setConnected(true);
    }
  }, [playerName, roomName]);

  if (!roomName || !playerName) {
    return <h1>404 - Invalid room name or player name</h1>;
  }

  return (
    <GameBoard 
      connected={connected}
      roomName={roomName}
      playerName={playerName}
    />
  )
}

interface GameBoardParam {
  connected: boolean;
  roomName: string;
  playerName: string;
}

function GameBoard(params: GameBoardParam) {
  const { connected, roomName, playerName } = params;

  const [playerNum, setPlayerNum] = useState(1);
  const [roomInfo, setRoomInfo] = useState<RoomInfoResponse>();
  const [turnInfo, setTurnInfo] = useState<TurnInfoMsgData>();

  // update message handlers
  useEffect(() => {
    if (!connected) {
      return;
    }

    const gameMsgHandler = gameService.getHandler<TopTenHandler>("game");
    const sysMsgHandler = gameService.getHandler<SystemMessageHandler>("system")
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
  }, [connected]);
  
  // get room info
  useEffect(() => {
    if (!connected) {
      return;
    }
    gameService.getRoomInfo(roomName)
      .then(info => {
        setPlayerNum(info.players.length);
        setRoomInfo(info);
      });
  }, [connected, roomName]);

  if (!connected) {
    return <h1>connecting</h1>;
  }

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