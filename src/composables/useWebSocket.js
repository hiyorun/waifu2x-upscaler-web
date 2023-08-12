import { ref } from "vue";

export function useWebSocket() {
  const socket = ref(null);

  const isConnected = ref(false);

  function handleConnection() {
    isConnected.value = true;
    console.log("Connected to WebSocket server.");
  }

  function handleDisconnection() {
    isConnected.value = false;
    console.log("Disconnected from WebSocket server.");
  }

  function handleError(event) {
    console.error("WebSocket Error:", event);
    handleDisconnection();
  }

  function createWebSocket() {
    socket.value = new WebSocket("ws://localhost:8080/api/v1/ws");
    socket.value.addEventListener("open", handleConnection);
    socket.value.addEventListener("close", handleDisconnection);
    socket.value.addEventListener("error", handleError);
  }

  function sendHeartbeat() {
    if (socket.value.readyState === WebSocket.OPEN) {
      socket.value.send("heartbeat");
    } else {
      console.log("Socket is not open for sending heartbeats.");
      createWebSocket();
    }
  }

  return {
    socket,
    handleConnection,
    handleDisconnection,
    handleError,
    createWebSocket,
    sendHeartbeat,
  };
}
