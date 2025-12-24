import { useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";
import { api } from "../api/api";
import { dtoToMessage, type MsgDto } from "../wire/message";


export default function Chat() {
  const { convs, msgs, loadMsgs, activeConvId, setActiveConvId } = useChat();
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    // if no convo selected, open the most recent one by default
    if (!activeConvId && convs.order.length) setActiveConvId(convs.order[0]);
  }, [activeConvId, convs.order]);
  
  const loadMore = async () => {
    if (!activeConvId || !hasMore) return;

    const currentMsgs = msgs[activeConvId];
    if (!currentMsgs || currentMsgs.order.length === 0) return;

    const oldestSeq = currentMsgs.order[0]; // Get the smallest sequence number
    
    // Fetch older messages
    try {
      const msgDto = await api<MsgDto[]>(`/message/before/${activeConvId}?limit=50&seq=${oldestSeq}`);
      // If fewer messages than the limit (50), we reached the beginning
      if (msgDto.length < 50) {
        setHasMore(false);
      }

      if (msgDto.length > 0) {
        const newMsgs = msgDto.map((w) => dtoToMessage(w));
        loadMsgs(newMsgs);
      }
    } catch (e) {
      console.error(e);
    }
  } 

  return (
    <div className="grid md:grid-cols-[320px_1fr] h-full overflow-hidden">
      <ConversationsSidebar
        convs={convs}
        activeConvId={activeConvId ?? undefined}
        onOpen={(cid) => setActiveConvId(cid)} // <-- no navigate
      />
      <section className="bg-white p-6 h-full overflow-hidden">
        {activeConvId ? (
          <MessagesPane
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
          />
        ) : (
          <div className="h-full grid place-items-center text-slate-400">
            Select a converseation to start chatting.
          </div>
        )}
      </section>
    </div>
  );
}
