const socket = new WebSocket("ws://localhost:8080/api/websocket/ws");
socket.onopen = () => {
  console.log("Connected to server ✅");
  socket.send("Hello from client!");
};

socket.onmessage = (e) => {
  console.log("Received:", e.data);
};

socket.onerror = (err) => {
  console.error("WebSocket error:", err);
};

socket.onclose = () => {
  console.log("Connection closed");
};