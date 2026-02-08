export type Message = {
  convId: string;
  fromUid: string;
  clientMsgId?: string; // Not presersisted on server, only for local pending messages
  ts: number; // used to order message when pending, updated from with server timestamp after acked

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
  created_at: number;
  // ***** if different type for message wanted, add Type field here
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

// ***** if the message is not text, think about how to handle it
export const jsonToMessage = (json: any): Message => ({
  id: json.server_msg_id,
  convId: json.conversation_id,
  fromUid: json.from_uid,
  clientMsgId: json.client_msg_id || null,
  seq: json.seq,
  ts: json.ts,
  text: json.body?.text ?? "",
});