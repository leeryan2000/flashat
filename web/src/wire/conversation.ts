import type { Conversation } from "../context/ChatContext";

export type ConvDto = {
  id: string;
  type: "group" | "direct";
  title:  string;
  avatar_url?: string | null;
  last_msg_id?: string;
  last_msg_text?: string;
  last_msg_from?: string;
  last_msg_ts?: number;
  last_seq: number;
  last_read_seq: number;
  unread_count: number;
};

export const toConversation = (w: ConvDto): Conversation => ({
  id: w.id,
  type: w.type,
  title: w.title,
  avatarUrl: w.avatar_url ?? null,
  lastMsgId: w.last_msg_id,
  lastMsgText: w.last_msg_text,
  lastMsgFrom: w.last_msg_from,
  lastMsgTs: w.last_msg_ts,
  lastSeq: w.last_seq,
  lastReadSeq: w.last_read_seq,
  unreadCount: w.unread_count,
});