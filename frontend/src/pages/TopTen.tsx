import { useEffect, useState } from "react"
import type { TopTenHandler, TopTenSender, GameInfo } from "../core/games/top10";
import type { RoomInfo, SystemMessageHandler } from "../core";
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
          gameService.joinGame(roomName, playerName)
            .then(() => {
              setConnected(true);
              window.addEventListener("beforeunload", gameService.leaveGame);
            }).catch((e) => alert(e));
          } else {
            setConnected(true);
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
        }).catch((e) => alert(e));
    } else {
      setConnected(true);
    }
  }, [playerName, roomName]);

  if (!roomName || !playerName) {
    return <h1>404 - Invalid room name or player name</h1>;
  }

  return (
    <>
      <nav>
        <div>room: {roomName}</div>
        <div>player: {playerName}</div>
      </nav>
      <GameBoard 
        connected={connected}
        playerName={playerName}
      />
    </>
  )
}

interface GameBoardParam {
  connected: boolean;
  playerName: string;
}

function GameBoard(params: GameBoardParam) {
  const { connected, playerName } = params;

  const [log, setLog] = useState("");
  const [inGame, setInGame] = useState(false);
  const [finished, setFinished] = useState(false);
  const [roomInfo, setRoomInfo] = useState<RoomInfo>();
  const [gameInfo, setGameInfo] = useState<GameInfo>();

  // update message handlers
  useEffect(() => {
    if (!connected) {
      return;
    }

    const sysMsgHandler = gameService.getHandler<SystemMessageHandler>("system")

    sysMsgHandler.onJoined((msg) => {
      setRoomInfo(msg.msg.roomInfo);
      setLog(`Player "${msg.msg.playerName}" has joined`);
    });
    sysMsgHandler.onLeft((msg) => {
      setRoomInfo(msg.msg.roomInfo);
      setLog(`Player "${msg.msg.playerName}" has joined`);
    });
    sysMsgHandler.onReady((msg) => {
      const newRoomInfo = msg.msg.roomInfo;
      setRoomInfo(newRoomInfo);
      setLog(`Player "${msg.msg.playerName}" ready for game (${newRoomInfo.players.length}/${newRoomInfo.roomSize})`);
    });
  }, [connected]);

  // update game message handlers
  useEffect(() => {
    if (!connected) {
      return;
    }

    const gameMsgHandler = gameService.getHandler<TopTenHandler>("game");

    gameMsgHandler.onGameInfo((msg) => {
      setGameInfo(msg.msg);
    });
    gameMsgHandler.onStart((msg) => {
      setGameInfo(msg.msg);
      setInGame(true);
      setLog("Game starts")
    });
    gameMsgHandler.onStartGuessing((msg) => {
      setGameInfo(msg.msg);
      setLog("Start guessing")
      console.log(msg.msg.numbers);
    });
    gameMsgHandler.onFinished(() => setFinished(true));
  }, [connected]);


  if (!connected) {
    return <h1>Not connected</h1>;
  }

  if (finished) {
    return <h1>Game finished</h1>;
  }

  const isGuesser: boolean = gameInfo?.guesser === playerName;

  return (
    <div>
      <h1>Top10</h1>

      {gameInfo? <h3>Turn {gameInfo.turn}/{gameInfo.maxTurn}</h3> : <></>}

      {isGuesser? 
        <Questions questions={gameInfo?.questions}/> : <></>
      }

      {gameInfo?.usedQuestion ? 
        <div>Question: {gameInfo?.usedQuestion}</div> : <></>
      }

      <div>Players: {roomInfo?.players.join(", ")}</div>

      {isGuesser?
        <h2>You are the guesser</h2> : <h2>Guesser: {gameInfo?.guesser}</h2>
      }
      
      {gameInfo?.numbers[playerName] ? <h2>Your Number: {gameInfo?.numbers[playerName]}</h2> : <></>}

      {inGame ? (
        <button onClick={() => {
          gameService.getSender<TopTenSender>().ready();
          setLog("Waiting for others to ready");
        }}>Ready</button>
      ) : (
        <button onClick={() => {
          gameService.ready();
          setLog("You are ready");
        }}>Start Game</button>
      )}

      <div>{log}</div>
    </div>
  );
}

interface QuestionsParam {
  questions?: string[];
}

function Questions(params: QuestionsParam) {
  const { questions } = params;
  if (!questions) {
    return <></>;
  }

  if (questions) {
    return (
      <div>
        {questions.map((question, index) => (
        <div key={`question-${index}`}>
          {question}
          <button 
            onClick={() => gameService.getSender<TopTenSender>().setQuestion(question)}>
              choose this question
          </button>
        </div>))}
      </div>
    );
  }

}