// const socket = new WebSocket("ws://localhost:8080/api/websocket/ws");
// socket.onopen = () => {
//   console.log("Connected to server ✅");
//   socket.send("Hello from client!");
// };

// socket.onmessage = (e) => {
//   console.log("Received:", e.data);
// };

// socket.onerror = (err) => {
//   console.error("WebSocket error:", err);
// };

// socket.onclose = () => {
//   console.log("Connection closed");
// };

const token1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhYTY2NmIxYS04YTViLTRlZDItYjAwNy1kZTdkOTc4MzJmMjciLCJpc3MiOiJGbGFzaGF0Iiwic3ViIjoidXNlci1hdXRoIiwiZXhwIjoxNzU0OTc4MzM1LCJpYXQiOjE3NTQ3MTkxMzV9.au2GpQfeGKvCV63p0aNBntKmuJFnEryvbfWaQ-rfeJ4"
const token2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI5NDc1ZjcwZC1kNzc0LTRlYzYtYTkzZi1lM2RmZDk1ZDVkZWEiLCJpc3MiOiJGbGFzaGF0Iiwic3ViIjoidXNlci1hdXRoIiwiZXhwIjoxNzU0OTg3MzExLCJpYXQiOjE3NTQ3MjgxMTF9.V6Vsi-cOBcIRE7oRNbrW8zMAM7Pn7GRDyPbqleoIWk8"
const ws1 = new WebSocket("ws://localhost:8080/api/websocket/ws", {
  headers: {
    Authorization: `Bearer ${token1}`,
  },
});

const ws2 = new WebSocket("ws://localhost:8080/api/websocket/ws", {
  headers: {
    Authorization: `Bearer ${token2}`,
  },
});

const queue = [];
let isOpen = true;

function safeSend(obj) {
  const payload = JSON.stringify(obj);
  if (isOpen) ws1.send(payload);
  else queue.push(payload);
}

function flush() {
  while (queue.length) ws1.send(queue.shift());
}


// High-level helpers
function joinRoom(roomId) {
  safeSend({ type: "join", roomId });
}

function sendChat(roomId, body) {
  safeSend({ type: "chat", roomId, body });
}

function sendDirect(toId, body) {
  safeSend({ type: "direct", toId, body });
}

ws1.addEventListener("open", () => {
  console.log("✅ Connected");
  joinRoom("general");
  sendChat("general", "Hello everyone! 👋");
  sendDirect("9475f70d-d774-4ec6-a93f-e3dfd95d5dea", "Hello user 9475f70d-d774-4ec6-a93f-e3dfd95d5dea! 👋");
});

ws1.addEventListener("message", (ev) => {
  try {
    const env = JSON.parse(ev.data);
    console.log("📩", env);
  } catch (e) {
    console.warn("⚠️ Bad message:", ev.data);
  }
});

ws2.addEventListener("message", (ev) => {
  try {
    const env = JSON.parse(ev.data);
    console.log("📩 user one received message", env);
  } catch (e) {
    console.warn("⚠️ Bad message:", ev.data);
  }
});

ws1.addEventListener("error", (err) => {
  console.error("❌ WebSocket error event:", err);
});

