export type Message = {
  convId: string;
  fromUid: string;
  fromName?: string;
  clientMsgId?: string; // Not persisted on server, only for local pending messages
  ts: number; // used to order message when pending, updated from server timestamp after acked

  id?: string; // server message id
  seq?: number; // server sequence number

  text?: string;
  status?: "sending" | "failed" | "sent"; // local only
};

export type MsgDto = {
  id: string;
  conversation_id: string;
  seq: number;
  from_uid: string;
  from_name: string;
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
  fromName: w.from_name,
  ts: w.created_at,
  text: w.body?.text ?? "",
});

export const jsonToMessage = (json: any): Message => ({
  id: json.server_msg_id,
  convId: json.conversation_id,
  fromUid: json.from_uid,
  fromName: json.from_name,
  clientMsgId: json.client_msg_id || null,
  seq: json.seq,
  ts: json.ts,
  text: json.body?.text ?? "",
});