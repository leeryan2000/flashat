import { useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";


export default function Chat() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { convs, msgs } = useChat();

  useEffect(() => {
    // if no convo selected, open the most recent one by default
    if (!selectedId && convs.order.length) setSelectedId(convs.order[0]);
  }, [selectedId, convs.order]);

  return (
    <div className="grid md:grid-cols-[320px_1fr] min-h-screen">
      <ConversationsSidebar
        convs={convs}
        activeConvId={selectedId ?? undefined}
        onOpen={(cid) => setSelectedId(cid)} // <-- no navigate
      />
      <section className="bg-white p-6">
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
