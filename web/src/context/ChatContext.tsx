import {
  createContext,
  use,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useAuth } from "./AuthContext";
import { api } from "../api/api";
import { toConversation, type ConversationWire } from "../wire/conversation";
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

type Message = {
  clientMsgId: string;
  convId: string;
  text: string;
  fromUid: string;
  id?: string; // server message id
  seq?: number; // would receive the seq from server
  ts?: number;
  status: "sending" | "failed" | "sent"; // local only
};

export type MsgSlice = {
  order: number[]; // ordered message ids
  entities: Record<number, Message>;
  pendingOrder: string[]; // ordered client message ids
  pendingEntities: Record<string, Message>;
};

export type MsgState = Record<string, MsgSlice>; // convId -> MsgSlice

type SendMsgInput = {
  convId: string;
  text: string;
};

interface ChatContext {
  convs: ConvState;
  msgs: Record<string, MsgSlice>; // convId -> MsgSlice

  loadConvs: (convs: Conversation[]) => void;
  // sendMessage: (input: SendMessageInput) => Promise<string>;
  // loadMessages: (messages: Message[]) => void; // sets up the messages when app started
  // markRead: (convId: string, seq: number) => void;
}

const ChatContext = createContext<ChatContext | undefined>(undefined);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [convs, setConvs] = useState<ConvState>({ entities: {}, order: [] });
  const [msgs, setMsgs] = useState<MsgState>({});

  // ***** should load the summary together with the connection of websocket
  useEffect(() => {
    // ***** Retrieve sidebar test
    const fetchUser = async () => {
      try {
        const summaryWire = await api<ConversationWire[]>(
          "/conversation/summary"
        );
        const summary = summaryWire.map(toConversation);
        loadConvs(summary);
      } catch (error) {
        console.error("Error fetching conversations:", error);
      } finally {
      }
    };
    fetchUser();
  }, []);

  useEffect(() => {
  }, [user]);

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
        // acked message from server 
        newState[convId].entities[msg.seq ?? 0] = msg;
        newState[convId].order.push(msg.seq ?? 0);
      }
      return newState;
    });
  }, []);


  const value = useMemo(
    () => ({
      convs,
      msgs,
      loadConvs,
      loadMsgs,
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
