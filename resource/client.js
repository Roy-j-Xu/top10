const host = window.location.hostname;
const socket = new WebSocket(`ws://${host}:8080/ws`)

const log = document.getElementById("log");
const playerInfoBoard = document.getElementById("player-info")
const guesserInfoBoard = document.getElementById("guesser")
const guesseNumberBoard = document.getElementById("guess-number")
const readyBtn = document.getElementById("readyBtn")
const questionBoard = document.getElementById("questions")

let ID = null
let GuessNumber = null

const MsgTypes = Object.freeze({
  JOINED: "joined",
  BROADCAST: "broadcast",
  QUESTOINS: "questions",
  GUESSER: "guesser",
  ASSIGN_NUMBER: "assign-number"
})

socket.onopen = () => {
  console.log("Connected")
  log.innerText = "Connected"
}

socket.onmessage = (event) => {
  console.log("Message from server:", event.data);

  const msgObj = JSON.parse(event.data);
  console.log(msgObj)
  let msg = msgObj["Msg"]

  switch (msgObj["Type"]) {
    case MsgTypes.JOINED:
      ID = msg
      playerInfoBoard.innerText = `Your ID: ${ID}`
      break
    case MsgTypes.QUESTOINS:
      if (msg["Guesser"] === ID) {
        playerInfoBoard.innerText = `Your ID: ${ID}. You are the guesser`
        questionBoard.innerText = msg["Questions"].join("\n\n")
      } else {
        playerInfoBoard.innerText = `Your ID: ${ID}. Player ${msg["Guesser"]} is the guesser`
        questionBoard.innerText = ""
      }
      break
    case MsgTypes.BROADCAST:
      log.innerText = msg + "\n" + log.innerText
      break
    case MsgTypes.ASSIGN_NUMBER:
      guesseNumberBoard.innerText = msg
      break
  }

}

function sendReady() {
  console.log("Sending message")
  socket.send(JSON.stringify({ type: "ready" }));
}

readyBtn.onclick = sendReady