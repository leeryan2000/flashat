import type { ConvState, MsgSlice, MsgState } from "../context/ChatContext";
import type { Conversation } from "../wire/conversation";
import type { Message } from "../wire/message";

export const applyAckedMsgToMsgState = (
    prev: MsgState,
    conversation_id: string,
    client_msg_id: string,
    seq: number,
    id: string
) => {
    const prevSlice = prev[conversation_id];
    if (!prevSlice) return prev; // no such conversation

    const msgToPromote = prevSlice.pendingEntities[client_msg_id];
    if (!msgToPromote) return prev; // no such pending message

    // set the acked message with server seq and id
    const ackedMsg: Message = { ...msgToPromote, status: "sent", seq, id };
    // Remove from pending
    const { [client_msg_id]: removed, ...remainingPendingEntities } =
        prevSlice.pendingEntities;

    const newSlice: MsgSlice = {
        // Only add seq to order if broadcast hasn't already added it
        order: prevSlice.order.includes(ackedMsg.seq!)
            ? prevSlice.order
            : [...prevSlice.order, ackedMsg.seq!],
        entities: { ...prevSlice.entities, [ackedMsg.seq!]: ackedMsg },
        pendingOrder: prevSlice.pendingOrder.filter(
            (cid) => cid !== client_msg_id
        ),
        pendingEntities: remainingPendingEntities,
    };

    return { ...prev, [conversation_id]: newSlice };
}


export const applyMsgToConvState = (
  prev: ConvState,
  msg: Message,
  currentUserId: string,
  activeConvId: string
) => {
  const conv = prev.entities[msg.convId];

  // Safety: If conversation doesn't exist yet (e.g., new chat), return state as is
  // (Or you could add logic here to create it on the fly)
  if (!conv) return prev;

  const isMyMsg = msg.fromUid === currentUserId;
  const isSelected = activeConvId === msg.convId;

  // Logic: Only increment unread if it's NOT my message AND I'm NOT looking at it
  const newUnread =
    !isMyMsg && !isSelected ? (conv.unreadCount || 0) + 1 : conv.unreadCount;
    
  const updatedConv = {
    ...conv,
    lastMsgText: msg.text,
    lastMsgTs: msg.ts,
    lastMsgFrom: msg.fromUid,
    unreadCount: newUnread,
    lastSeq: msg.seq ?? conv.lastSeq,
  };

  const newEntities = { ...prev.entities, [msg.convId]: updatedConv };

  // Re-sort: Move this conversation to the top
  // Assuming getSortedConvIds is your existing sorting helper
  const newOrder = getSortedConvIds(newEntities);

  return { entities: newEntities, order: newOrder };
};

export function getSortedConvIds(
  entities: Record<string, Conversation>
): string[] {
  return Object.values(entities)
    .sort((a, b) => {
      const diff = (b.lastMsgTs ?? 0) - (a.lastMsgTs ?? 0);
      return diff !== 0 ? diff : a.id.localeCompare(b.id);
    })
    .map((c) => c.id);
}
