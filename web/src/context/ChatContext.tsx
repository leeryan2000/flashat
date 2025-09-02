import { createContext, use, useCallback, useState } from "react";
import { useAuth } from "./AuthContext";
import { v4 } from "uuid";
// create types that match exactly what the server passed in

type Conversation = {
  id: string;
  title: string;
  lastSeq?: number;
  lastReadSeq?: number;
  unread?: number;
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
  conversations: Record<string, Conversation>;
  messages: Record<string, Message>;

  sendMessage: (input: SendMessageInput) => Promise<string>;
  loadConversations: (conversations: Conversation[]) => void;
  loadMessages: (messages: Message[]) => void; // sets up the messages when app started
  markRead: (convId: string, seq: number) => void;
}

const ChatContext = createContext<ChatContext | undefined>(undefined);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [conversations, setConversations] = useState<Record<string, Conversation>>({});
  const [messages, setMessages] = useState<Record<string, Message>>({});

  // ***** use websocket probably with another context
  const sendMessage = useCallback(
    async (input: SendMessageInput): Promise<string> => {
      const clientMsgId = v4();
      const fromUid = user?.uid || "unknown";
      const optimistic: Message = {
        clientMsgId,
        convId: input.convId,
        fromUid: user?.uid || "unknown",
        text: input.text,
        seq: conversations[input.convId]?.lastSeq || 0,
        status: "sending"
      };

      return "";
    },
    []
  );
}
