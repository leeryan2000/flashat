import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type PropsWithChildren,
} from "react";

type WSStatus = "idle" | "connecting" | "open" | "closed" | "error";

type WebSocketContextType = {
  status: WSStatus;
  lastMsg: string | null;
  send: (data: string) => void;
  disconnect: () => void;
};

const WebSocketContext = createContext<WebSocketContextType | null>(null);
type WebSocketProviderProps = PropsWithChildren<{
  url: string;
}>;

export function WebSocketProvider({ url, children }: WebSocketProviderProps) {
  const [status, setStatus] = useState<WSStatus>("idle");
  const [lastMsg, setLastMsg] = useState<string | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  // ***** implement auto reconnection to websocket when offlined
  useEffect(() => {
    const uuid = crypto.randomUUID();
    console.log(`Establishing WebSocket connection with UUID: ${uuid}`);
    
    const ws = new WebSocket(url);
    socketRef.current = ws;

    setStatus("connecting");

    ws.onopen = () => setStatus("open");
    ws.onclose = () => {
      setStatus("closed");
      console.log(`WS close id=`, uuid)
    }; 
    // Set lastMsg when a message is received
    ws.onmessage = (wire) => {
      setLastMsg(wire.data);
    };

    return () => {
      ws.close();
      socketRef.current = null;
      console.log("WebSocket closed on component unmount id:", uuid);
      setStatus("closed");
    };
  }, []);

  const send = useCallback((data: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(data);
      console.log("Sent message:", data);
    } else {
      console.warn("WebSocket not open; message dropped.");
    }
  }, []);

  const disconnect = useCallback(() => {
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
    [status, lastMsg]
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
