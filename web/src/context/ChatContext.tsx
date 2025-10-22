import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useAuth } from "./AuthContext";
import { api } from "../api/api";
import { toConversation, type ConversationWire } from "../wire/conversation";
import { toMessage, type MsgWire } from "../wire/message";
import { useWebSocket } from "./WebSocketContext";
// create types that match exactly what the server passed in

export type Conversation = {
  id: string;
  type: "group" | "direct";

  title: string;
  avatarUrl?: string | null;
  lastMsgId?: string;
  lastMsgText?: string;
  lastMsgFrom?: string;
  // the timestamp in database are timestampz match it
  lastMsgTs?: number;
  lastSeq: number;

  lastReadSeq: number;
  unreadCount: number;
};

export type ConvState = {
  order: string[]; // the order of conversation ids, sorted by timestamp of last message
  entities: Record<string, Conversation>;
};

export type Message = {
  convId: string;
  fromUid: string;
  clientMsgId?: string;
  ts: number; // used to order message when pending, updated from with server timestamp after acked
  
  id?: string; // server message id
  seq?: number; // server sequence number
  
  text?: string;
  status?: "sending" | "failed" | "sent"; // local only
};

export type MsgSlice = {
  order: number[]; // ordered message seq
  entities: Record<number, Message>;
  pendingOrder: string[]; // ordered client message ids
  pendingEntities: Record<string, Message>;
};

export type MsgState = Record<string, MsgSlice>; // convId -> MsgSlice

interface ChatContext {
  convs: ConvState;
  msgs: Record<string, MsgSlice>; // convId -> MsgSlice

  loadConvs: (convs: Conversation[]) => void;
  loadMsgs: (messages: Message[]) => void;
  sendMessage: (text: string, convId: string) => void;
}

const ChatContext = createContext<ChatContext | null>(null);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const { send } = useWebSocket();
  const { lastMsg } = useWebSocket();
  const [convs, setConvs] = useState<ConvState>({ entities: {}, order: [] });
  const [msgs, setMsgs] = useState<MsgState>({});

  useEffect(() => {
    const controller = new AbortController();
    const { signal } = controller;

    (async () => {
      try {
        const summaryWire = await api<ConversationWire[]>(
          "/conversation/summary",
          { signal } as any // to allow the controller.abort to send signal
        );
        const summary = summaryWire.map(toConversation);
        loadConvs(summary);
      } catch (err) {
        console.error("Error fetching conversations:", err);
      }
    })();

    return () => controller.abort();
  }, []);
  
  // ***** modify
  useEffect(() => {
    const controller = new AbortController();
    const { signal } = controller;
    
    (async () => {
      if (convs.order.length > 0) {
        try {
          const msgWire = await api<MsgWire[]>(
            `/message/latest/${convs.order[0]}?limit=50`,
            { signal } as any
          );
          const msgs = msgWire.map((w) => toMessage(w));
          console.log("Fetched messages:", msgs);
          loadMsgs(msgs);

        } catch (err) {
          console.error("Error fetching messages:", err);
        }
      }
    })();

    return () => controller.abort();
  }, [convs.order]);

  useEffect(() => {
    if (!lastMsg) return;
    const jsonMsg = JSON.parse(lastMsg);
    


  }), [lastMsg];

  const loadConvs = useCallback((list: Conversation[]) => {
    setConvs((prev) => {
      const entities = { ...prev.entities };
      // overwrite the conversation state with the new data
      for (const conv of list) {
        entities[conv.id] = { ...(entities[conv.id] ?? {}), ...conv };
      }

      const order = Object.values(entities)
        .sort((a, b) => {
          const diff = (b.lastMsgTs ?? 0) - (a.lastMsgTs ?? 0);
          return diff !== 0 ? diff : a.id.localeCompare(b.id);
        })
        .map((conv) => conv.id);
      return {
        entities: entities,
        order: order,
      };
    });
  }, []);

  const loadMsgs = useCallback((list: Message[]) => {
    setMsgs((prev) => {
      const newState = { ...prev };
      for (const msg of list) {
        const convId = msg.convId;
        if (!newState[convId]) {
          newState[convId] = {
            order: [],
            entities: {},
            pendingOrder: [],
            pendingEntities: {},
          };
        }

        msg.status = "sent";
        
        // acked message from server 
        if (msg.seq && !newState[convId].entities[msg.seq]) {
          // update the existing message with id and seq
          newState[convId].order.push(msg.seq ?? 0);
          newState[convId].entities[msg.seq ?? 0] = msg;
        } 
      }
      return newState;
    });
  }, []);

  const sendMessage = useCallback((text: string, convId: string) => {
    if (text.trim() === "") return;
    const message = { 
      type: "chat",
      conversation_id: convId,
      client_msg_id: crypto.randomUUID(),
      from_uid: user?.uid,
      ts: Date.now(),
      body: {
        text: text.trim()
      }
    };
    send(JSON.stringify(message));
  }, [send, user]);

  // Websocket Connection

  const value = useMemo(
    () => ({
      convs,
      msgs,
      loadConvs,
      loadMsgs,
      sendMessage,
    }),
    [convs, msgs]
  );

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
}

export function useChat() {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error("useChat must be used within a ChatProvider");
  }
  return context;
}

