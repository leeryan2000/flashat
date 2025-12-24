import type { Message } from "../context/ChatContext";

export type MsgDto = {
  id: string;
  conversation_id: string;
  seq: number;
  from_uid: string;
  created_at: number;
  body: {
    text: string;
  };
};

export const dtoToMessage = (w: MsgDto): Message => ({
  id: w.id,
  convId: w.conversation_id,
  seq: w.seq,
  fromUid: w.from_uid,
  ts: w.created_at,
  text: w.body?.text ?? "", // safely get text from body
});

export const jsonToMessage = (json: any): Message => ({
  id: json.server_msg_id,
  convId: json.conversation_id,
  seq: json.seq,
  fromUid: json.from_uid,
  ts: json.ts,
  text: json.body?.text ?? "",
});