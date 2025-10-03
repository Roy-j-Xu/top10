const host = window.location.hostname;
const socket = new WebSocket(`ws://${host}:8080/ws`)

const log = document.getElementById("log");
const readyBtn = document.getElementById("readyBtn")
const questionBoard = document.getElementById("questions")

socket.onopen = () => {
  console.log("Connected")
  log.innerText = "Connected"
}

socket.onmessage = (event) => {
  console.log("Message from server:", event.data);

  const msgObj = JSON.parse(event.data);

  switch (msgObj["Type"]) {
    case "broadcast":
      break
    case "questions":
  }

  log.innerText = msgObj["Msg"];
}

function sendReady() {
  console.log("Sending message")
  socket.send(JSON.stringify({ type: "ready" }));
}

readyBtn.onclick = sendReady