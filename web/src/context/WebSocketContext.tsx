import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

type WSStatus = "idle" | "connecting" | "open" | "closed" | "error";

type WebSocketContextType = {
  status: WSStatus;
  lastMsg: string | null;
  send: (data: string) => void;
  disconnect: () => void;
};

const WebSocketContext = createContext<WebSocketContextType | null>(null);

const getWebSocketURL = () => {
  const envUrl = import.meta.env.VITE_WS_URL;

  // 1. Use the .env value if it exists (for local development)
  if (envUrl) return envUrl;

  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const host = window.location.host; 
  
  return `${protocol}//${host}/api/websocket/ws`;
};

export function WebSocketProvider({ children }: {children: React.ReactNode}) {
  const [status, setStatus] = useState<WSStatus>("idle");
  const [lastMsg, setLastMsg] = useState<string | null>(null);
  
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);
  
  // Track if the disconnection was intentional (user clicked logout/disconnect)
  const isIntentionalClose = useRef(false);

  // Configuration
  const RECONNECT_INTERVAL = 5000; // ms before attempting reconnect

  const connect = useCallback(() => {
    // Prevent multiple connections
    if (socketRef.current?.readyState === WebSocket.OPEN) return;

    // Clear any existing reconnect timers
    if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);

    const uuid = crypto.randomUUID();
    console.log(`Attempting WebSocket connection... ID: ${uuid}`);
    
    const ws = new WebSocket(getWebSocketURL());
    socketRef.current = ws;
    setStatus("connecting");
    isIntentionalClose.current = false;

    ws.onopen = () => {
      console.log("WebSocket Connected");
      setStatus("open");
    };

    ws.onmessage = (wire) => {
      setLastMsg(wire.data);
    };

    ws.onclose = (event) => {
      console.log(`WebSocket closed. Code: ${event.code}, Reason: ${event.reason}`);
      setStatus("closed");
      
      // Only reconnect if it wasn't an intentional disconnect
      if (!isIntentionalClose.current) {
        console.log(`Reconnecting in ${RECONNECT_INTERVAL}ms...`);
        reconnectTimeoutRef.current = window.setTimeout(() => {
          connect();
        }, RECONNECT_INTERVAL);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

  }, [getWebSocketURL]);

  // Initial Connection
  useEffect(() => {
    connect();

    // Cleanup on unmount
    return () => {
      isIntentionalClose.current = true;
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
      socketRef.current?.close();
    };
  }, [connect]);

  const send = useCallback((data: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(data);
    } else {
      console.warn("WebSocket not open; message dropped.");
    }
  }, []);

  const disconnect = useCallback(() => {
    isIntentionalClose.current = true; // Prevent auto-reconnect
    if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
    
    socketRef.current?.close();
    setStatus("closed");
  }, []);

  const value = useMemo(
    () => ({
      status,
      lastMsg,
      send,
      disconnect,
    }),
    [status, lastMsg, send, disconnect]
  );

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context)
    throw new Error("useWebSocket must be used within WebSocketProvider");
  return context;
}
