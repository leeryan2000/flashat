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
  entities: Record<string, Conversation>;
  order: string[];
};

type Message = {
  id?: string; // server message id
  clientMsgId?: string;
  convId: string;
  text: string;
  fromUid: string;
  ts?: number;
  seq: number; // would receive the seq from server
  status?: "sending" | "failed" | "sent"; // local only
};

type SendMessageInput = {
  convId: string;
  text: string;
};

interface ChatContext {
  convs: ConvState;
  msgs: Record<string, Message>;

  loadConvs: (convs: Conversation[]) => void;
  // sendMessage: (input: SendMessageInput) => Promise<string>;
  // loadMessages: (messages: Message[]) => void; // sets up the messages when app started
  // markRead: (convId: string, seq: number) => void;
}

const ChatContext = createContext<ChatContext | undefined>(undefined);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [convs, setConvs] = useState<ConvState>({ entities: {}, order: [] });
  const [msgs, setMessages] = useState<Record<string, Message>>({});

  useEffect(() => {
    // ***** Retrieve sidebar test
    const fetchUser = async () => {
      try {
        const summaryWire = await api<ConversationWire[]>("/conversation/summary");
        const summary = summaryWire.map(toConversation);
        loadConvs(summary);
      } catch (error) {
        console.error("Error fetching conversations:", error);
      } finally {
      }
    };
    fetchUser();
  }, []);

  const loadConvs = useCallback((list: Conversation[]) => {
    setConvs((prev) => {
      const convs = { ...prev.entities };
      for (const conv of list) {
        convs[conv.id] = { ...(convs[conv.id] ?? {}), ...conv };
      }
      const order = Object.values(convs)
        .sort((a, b) => {
          const diff = (b.lastMsgTs ?? 0) - (a.lastMsgTs ?? 0);
          return diff !== 0 ? diff : a.id.localeCompare(b.id);
        })
        .map((conv) => conv.id);
      return {
        entities: convs,
        order: order,
      };
    });
  }, []);

  const value = useMemo(
    () => ({
      convs,
      msgs,
      loadConvs,
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
