import { useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";
import { api } from "../api/api";
import { dtoToMessage, type MsgDto } from "../wire/message";
import { useFriendshipActions } from "../hooks/useFriendshipActions";
import type { CustomResponse } from "../wire/resp";

export default function Chat() {
  const { convs, msgs, loadMsgs, activeConvId, setActiveConvId, friendships, removeConv } = useChat();
  const { blockUser } = useFriendshipActions();
  const [hasMore, setHasMore] = useState(true);

  // reset
  useEffect(() => {
    setHasMore(true);
  }, [activeConvId]);

  const activeConv = activeConvId ? convs.entities[activeConvId] ?? null : null;

  // For direct convs, find the other participant via the friendships list
  const activeFriend = activeConvId
    ? friendships.find((f) => f.directConversationId === activeConvId) ?? null
    : null;

  const handleBlock = async () => {
    if (!activeFriend || !activeConvId) return;
    await blockUser(activeFriend.uid, activeConvId);
    setActiveConvId(null);
  };

  const handleLeaveGroup = async () => {
    if (!activeConvId) return;
    try {
      await api<CustomResponse>(`/conversation/${activeConvId}/leave`, { method: "DELETE" });
      removeConv(activeConvId);
      setActiveConvId(null);
    } catch (err) {
      console.error("Error leaving group:", err);
    }
  };

  const loadMore = async () => {
    if (!activeConvId || !hasMore) return;

    const currentMsgs = msgs[activeConvId];
    if (!currentMsgs || currentMsgs.order.length === 0) return;

    const oldestSeq = currentMsgs.order[0];

    try {
      const msgDto = await api<MsgDto[]>(`/message/before/${activeConvId}?limit=50&seq=${oldestSeq}`);
      if (msgDto.length < 50) {
        setHasMore(false);
      }
      if (msgDto.length > 0) {
        loadMsgs(msgDto.map(dtoToMessage));
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="grid md:grid-cols-[320px_1fr] h-full overflow-hidden">
      <ConversationsSidebar
        convs={convs}
        activeConvId={activeConvId ?? undefined}
        onOpen={(cid) => setActiveConvId(cid)}
      />
      <section className="bg-white h-full overflow-hidden">
        {activeConvId ? (
          <MessagesPane
            key={activeConvId}
            conv={activeConv ?? undefined}
            msg={
              msgs[activeConvId] ?? {
                order: [],
                entities: {},
                pendingOrder: [],
                pendingEntities: {},
              }
            }
            activeConvId={activeConvId}
            onLoadMore={loadMore}
            onBlock={activeFriend ? handleBlock : undefined}
            onLeaveGroup={activeConv?.type === "group" ? handleLeaveGroup : undefined}
          />
        ) : (
          <div className="h-full grid place-items-center text-slate-400">
            Select a conversation to start chatting.
          </div>
        )}
      </section>
    </div>
  );
}
