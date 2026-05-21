import { useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";
import { api } from "../api/api";
import { dtoToMessage, type MsgDto } from "../wire/message";
import { useFriendshipActions } from "../hooks/useFriendshipActions";
import type { CustomResponse } from "../wire/resp";
import { useAuth } from "../context/AuthContext";

export default function Chat() {
  const { convs, msgs, loadMsgs, activeConvId, setActiveConvId, friendships, removeConv } = useChat();
  const { blockUser } = useFriendshipActions();
  const { user } = useAuth();
  const [hasMore, setHasMore] = useState(true);
  const [leaveConfirm, setLeaveConfirm] = useState(false);

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

  const handleLeaveGroup = () => {
    if (!activeConvId) return;
    setLeaveConfirm(true);
  };

  const confirmLeave = async () => {
    if (!activeConvId) return;
    setLeaveConfirm(false);
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

  const isLastMember = (activeConv?.participants.filter(p => p.uid !== user?.uid).length ?? 1) === 0;

  return (
    <div className="grid md:grid-cols-[320px_1fr] h-full overflow-hidden">
      {/* Leave group confirmation dialog */}
      {leaveConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-80 flex flex-col gap-4">
            <h2 className="font-semibold text-slate-800 text-base">Leave Group</h2>
            {isLastMember ? (
              <p className="text-sm text-slate-600">
                You are the last member. Leaving will permanently delete this conversation and all its messages.
              </p>
            ) : (
              <p className="text-sm text-slate-600">
                Are you sure you want to leave this group?
              </p>
            )}
            <div className="flex gap-2 justify-end">
              <button
                onClick={() => setLeaveConfirm(false)}
                className="px-4 py-2 rounded-xl text-sm text-slate-600 hover:bg-slate-100 transition"
              >
                Cancel
              </button>
              <button
                onClick={confirmLeave}
                className="px-4 py-2 rounded-xl text-sm bg-red-500 text-white hover:bg-red-600 transition"
              >
                {isLastMember ? "Delete & Leave" : "Leave"}
              </button>
            </div>
          </div>
        </div>
      )}
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
