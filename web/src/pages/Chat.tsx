import { useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";
import { api } from "../api/api";
import { toMessage, type MsgDto } from "../wire/message";


export default function Chat() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { convs, msgs, loadMsgs} = useChat();

  useEffect(() => {
    // if no convo selected, open the most recent one by default
    if (!selectedId && convs.order.length) setSelectedId(convs.order[0]);
  }, [selectedId, convs.order]);
  
  const loadMore = async () => {
    if (!selectedId) return;
    const currentMsgs = msgs[selectedId];
    if (!currentMsgs || currentMsgs.order.length === 0) return;

    const oldestSeq = currentMsgs.order[0]; // Get the smallest sequence number
    
    // Fetch messages OLDER than the oldest one we have
    // Assuming your API supports ?beforeSeq=...
    const msgDto = await api<MsgDto[]>(`/message/before/${selectedId}?limit=50&seq=${oldestSeq}`);
    const newMsgs = msgDto.map((w) => toMessage(w));
    console.log("Fetched messages:", newMsgs);
    loadMsgs(newMsgs);
  } 

  return (
    <div className="grid md:grid-cols-[320px_1fr] h-full overflow-hidden">
      <ConversationsSidebar
        convs={convs}
        activeConvId={selectedId ?? undefined}
        onOpen={(cid) => setSelectedId(cid)} // <-- no navigate
      />
      <section className="bg-white p-6 h-full overflow-hidden">
        {selectedId ? (
          <MessagesPane
            msg={
              msgs[selectedId] ?? {
                order: [],
                entities: {},
                pendingOrder: [],
                pendingEntities: {},
              }
            }
            activeConvId={selectedId}
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
