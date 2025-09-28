import type { Message } from "../context/ChatContext";

export type MessageWire = {
  id: string;
  conversation_id: string;
  seq: number;
  from_uid: string;
  body: {
    text: string;
  }
  created_at: number;
};

export const toMessage = (w: MessageWire) : Message => ({
    convId: w.conversation_id,
    id: w.id,
    fromUid: w.from_uid,
    text: w.body?.text ?? "", // safely get text from body
    ts: w.created_at,
    self: false,
});
