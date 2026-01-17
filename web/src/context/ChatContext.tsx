import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { useAuth } from "./AuthContext";
import { api } from "../api/api";
import { toConversation, type ConvDto, type Conversation } from "../wire/conversation";
import { dtoToMessage, jsonToMessage, type Message, type MsgDto } from "../wire/message";
import { useWebSocket } from "./WebSocketContext";
import { applyAckedMsgToMsgState, applyMsgToConvState, getSortedConvIds, } from "../utils/chatHelpers";
// create types that match exactly what the server passed in

export type ConvState = {
  order: string[]; // the order of conversation ids, sorted by timestamp of last message
  entities: Record<string, Conversation>;
};

// Messages stored by conversation id
export type MsgState = Record<string, MsgSlice>; // convId -> MsgSlice

export type MsgSlice = {
  order: number[]; // ordered message seq
  entities: Record<number, Message>;
  pendingOrder: string[]; // ordered client message ids
  pendingEntities: Record<string, Message>;
};

interface ChatContext {
  convs: ConvState;
  msgs: Record<string, MsgSlice>; // convId -> MsgSlice
  activeConvId: string | null;

  loadConvs: (convs: Conversation[]) => void;
  loadMsgs: (messages: Message[]) => void;
  sendMessage: (text: string, convId: string) => void;
  setActiveConvId(id: string | null): void;
  markAsRead(convId: string): void;
}

const ChatContext = createContext<ChatContext | null>(null);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const { send, lastMsg, status } = useWebSocket();
  const [convs, setConvs] = useState<ConvState>({ entities: {}, order: [] });
  const [msgs, setMsgs] = useState<MsgState>({}); // Messages for the client stored here
  const [activeConvId, setActiveConvId] = useState<string | null>(null);

  const hasInitialFetched = useRef(false);

  const markAsRead = useCallback(
    (convId: string) => {
      const conv = convs.entities[convId];

      if (conv && conv.unreadCount > 0) {
        // Optimistic Update
        setConvs((prev) => {
          const current = prev.entities[convId];
          if (!current) return prev;

          return {
            ...prev,
            entities: {
              ...prev.entities,
              [convId]: {
                ...current,
                unreadCount: 0,
                lastReadSeq: current.lastSeq,
              },
            },
          };
        });

        // Send to Server
        send(
          JSON.stringify({
            type: "read",
            conversation_id: convId,
            from_uid: user?.uid,
            last_read_seq: conv.lastSeq,
          })
        );
      }
    },
    [convs, send, user]
  );

  useEffect(() => {
    if (status === "closed") return;

    const controller = new AbortController();
    const { signal } = controller;

    (async () => {
      try {
        const summaryDto = await api<ConvDto[]>(
          "/conversation/summary",
          { signal } as any // to allow the controller.abort to send signal
        );
        const summary = summaryDto.map(toConversation);
        loadConvs(summary);
      } catch (err: any) {
        if (err.name !== "AbortError") {
          console.error("Error fetching conversations:", err);
        }
      }
    })();

    return () => controller.abort();
  }, [status]);

  // Fectch messages for the most recent conversation
  useEffect(() => {
    // if conversations aren't loaded yet, skip
    if (hasInitialFetched.current || convs.order.length === 0) return;

    hasInitialFetched.current = true;

    const controller = new AbortController();
    const { signal } = controller;

    (async () => {
      // Create an array of promises to fetch all concurrently
      const fetchPromises = convs.order.map(async (convId) => {
        try {
          const msgDto = await api<MsgDto[]>(
            `/message/latest/${convId}?limit=50`,
            { signal } as any
          );
          const msgs = msgDto.map((w) => dtoToMessage(w));

          loadMsgs(msgs);
        } catch (err: any) {
          // Ignore abort errors, log others
          if (err.name !== "AbortError") {
            console.error(`Error fetching messages for ${convId}:`, err);
          }
        }
      });

      // Wait for all to finish (optional, but good for cleanup logic if needed)
      await Promise.all(fetchPromises);
    })();

    return () => controller.abort();
  }, [convs.order]);

  // WebSocket message handler when lastMsg changes
  useEffect(() => {
    if (!lastMsg) return;
    const jsonMsg = JSON.parse(lastMsg);

    if (jsonMsg.type === "ack") {

      const { conversation_id, client_msg_id, seq, id } = jsonMsg;
      setMsgs((prev) => {
        return applyAckedMsgToMsgState(prev, conversation_id, client_msg_id, seq, id);
      });

    } else if (jsonMsg.type === "chat") {
      const msg = jsonToMessage(jsonMsg);
      loadMsgs([msg]);
      setConvs((prev) => {
        return applyMsgToConvState(prev, msg, user?.uid || "", activeConvId || "");
      });
    }
  }, [lastMsg, user]);

  

  const loadConvs = useCallback((list: Conversation[]) => {
    setConvs((prev) => {
      const entities = { ...prev.entities };
      // overwrite the conversation state with the new data
      for (const conv of list) {
        entities[conv.id] = { ...(entities[conv.id] ?? {}), ...conv };
      }

      // Sort the conversation everytime received based on the last message timestamp
      const order = getSortedConvIds(entities);
      return {
        entities: entities,
        order: order,
      };
    });
  }, []);

  const loadMsgs = useCallback((list: Message[]) => {
    if (list.length === 0) return;

    const convId = list[0].convId;

    setMsgs((prev) => {
      const prevSlice = prev[convId] ?? {
        order: [],
        entities: {},
        pendingOrder: [],
        pendingEntities: {},
      };
      const newOrder: number[] = [];
      const newEntities: Record<number, Message> = {};
      let hasUpdates = false;

      for (const incoming of list) {
        // if missing seq, skip
        if (!incoming.seq) {
          continue;
        }

        const existing = prevSlice.entities[incoming.seq];

        const msg: Message = {
          ...(existing ?? {}),
          ...incoming,
          status: "sent",
        };

        // push seq to order only if it's new
        if (!existing) {
          newOrder.push(incoming.seq);
        }

        newEntities[incoming.seq] = msg;
        hasUpdates = true;
      }

      // if no updates detected, return previous state
      if (!hasUpdates) return prev;

      // Sort combined order
      const nextOrder = [...prevSlice.order, ...newOrder].sort((a, b) => a - b);

      return {
        ...prev,
        [convId]: {
          ...prevSlice,
          order: nextOrder,
          entities: { ...prevSlice.entities, ...newEntities }, // Overwrite with new/updated entities
        },
      };
    });
  }, []);

  const sendMessage = useCallback(
    (text: string, convId: string) => {
      if (text.trim() === "") return;

      const clientMsgId = crypto.randomUUID();
      const localTs = Date.now();

      // Match server message type
      const message = {
        type: "chat",
        conversation_id: convId,
        client_msg_id: clientMsgId,
        from_uid: user?.uid,
        ts: localTs,
        body: {
          text: text.trim(),
        },
      };

      // Put the message in to the pending list
      setMsgs((prev) => {
        const prevSlice = prev[convId];
        if (!prevSlice) return prev;

        const newMsg: Message = {
          convId: convId,
          fromUid: user?.uid || "",
          clientMsgId: clientMsgId,
          ts: message.ts,
          text: text.trim(),
          status: "sending",
        };
        const newSlice: MsgSlice = {
          order: prevSlice.order,
          entities: prevSlice.entities,
          pendingOrder: [...prevSlice.pendingOrder, clientMsgId],
          pendingEntities: {
            ...prevSlice.pendingEntities,
            [clientMsgId]: newMsg,
          },
        };
        return { ...prev, [convId]: newSlice };
      });

      send(JSON.stringify(message));

      // Set a timeout to mark the message as 'failed' if no ACK received
      setTimeout(() => {
        setMsgs((prev) => {
          const slice = prev[convId];
          if (!slice) return prev;

          const pendingMsg = slice.pendingEntities[clientMsgId];

          if (!pendingMsg) return prev;
          return {
            ...prev,
            [convId]: {
              ...slice,
              pendingEntities: {
                ...slice.pendingEntities,
                [clientMsgId]: { ...pendingMsg, status: "failed" },
              },
            },
          };
        });
      }, 10000); // 10 seconds timeout
    },
    [send, user]
  );

  // Websocket Connection

  const value = useMemo(
    () => ({
      convs,
      msgs,
      activeConvId,
      loadConvs,
      loadMsgs,
      sendMessage,
      setActiveConvId,
      markAsRead,
    }),
    [
      convs,
      msgs,
      activeConvId,
      loadConvs,
      loadMsgs,
      sendMessage,
      setActiveConvId,
      markAsRead,
    ]
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

