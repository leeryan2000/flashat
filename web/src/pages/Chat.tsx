import { useMemo, useRef, useState, useEffect } from "react";
import ConversationsSidebar from "../components/ConversationSidebar";
import { useChat, type ConvState } from "../context/ChatContext";
import MessagesPane from "../components/MessagePane";

/**
 * Minimal, self‑contained chat page with:
 * - Left sidebar: conversation list (search + unread badge)
 * - Right pane: messages for the selected conversation + composer
 * - Routing: /chat and /chat/:id (useParams + navigate on row click)
 * - Tailwind for styling
 *
 * Drop this file into your project (e.g., src/pages/ChatPage.tsx),
 * add a route for it, and swap the mocked data for your API calls.
 */

// ---------- Types ----------
// export interface ConversationDemo {
//   id: string;
//   title: string; // display name / group name
//   last: string; // last message preview
//   ts: number; // unix ms
//   unread?: number; // unread count
// }

// export interface Message {
//   id: string;
//   convId: string;
//   from: string; // uid or display name
//   text: string;
//   ts: number; // unix ms
//   self?: boolean; // render bubble on the right
// }

// ---------- Mock data (replace with your API) ----------
// const MOCK_CONVOS: ConversationDemo[] = [
//   { id: "c1", title: "Alice", last: "See you tomorrow!", ts: Date.now() - 1000 * 60 * 3, unread: 2 },
//   { id: "c2", title: "Bob • Project", last: "Pushed a new branch.", ts: Date.now() - 1000 * 60 * 11 },
//   { id: "c3", title: "Team Rocket", last: "Prepare for trouble.", ts: Date.now() - 1000 * 60 * 60 * 5, unread: 7 },
//   { id: "c4", title: "Huai'en (You)", last: "Drafting the layout now", ts: Date.now() - 1000 * 60 * 60 * 10 },
// ];

// export const MOCK_CONV_STATE: ConvState = {
//   entities: {
//     c1: {
//       id: "c1",
//       type: "direct",
//       title: "Alice",
//       avatarUrl: "https://example.com/avatars/alice.png",
//       lastMsgId: "msg-c1-0120",
//       lastMsgText: "See you tomorrow!",
//       lastMsgFrom: "alice-uid",
//       lastMsgTs: Date.now() - 1000 * 60 * 3, // 3 minutes ago
//       lastSeq: 120,
//       lastReadSeq: 118,
//       unreadCount: 2,
//     },
//     c2: {
//       id: "c2",
//       type: "direct",
//       title: "Bob • Project",
//       avatarUrl: "https://example.com/avatars/bob.png",
//       lastMsgId: "msg-c2-0085",
//       lastMsgText: "Pushed a new branch.",
//       lastMsgFrom: "bob-uid",
//       lastMsgTs: Date.now() - 1000 * 60 * 11, // 11 minutes ago
//       lastSeq: 85,
//       lastReadSeq: 85,
//       unreadCount: 0,
//     },
//     c3: {
//       id: "c3",
//       type: "group",
//       title: "Team Rocket",
//       avatarUrl: "https://example.com/avatars/team-rocket.png",
//       lastMsgId: "msg-c3-0402",
//       lastMsgText: "Prepare for trouble.",
//       lastMsgFrom: "jessie-uid",
//       lastMsgTs: Date.now() - 1000 * 60 * 60 * 5, // 5 hours ago
//       lastSeq: 402,
//       lastReadSeq: 395,
//       unreadCount: 7,
//     },
//     c4: {
//       id: "c4",
//       type: "direct",
//       title: "Huai'en (You)",
//       avatarUrl: null,
//       lastMsgId: "msg-c4-0260",
//       lastMsgText: "Drafting the layout now",
//       lastMsgFrom: "you-uid",
//       lastMsgTs: Date.now() - 1000 * 60 * 60 * 10, // 10 hours ago
//       lastSeq: 260,
//       lastReadSeq: 260,
//       unreadCount: 0,
//     },
//   },
//   // Sorted by recency (most recent first)
//   order: ["c1", "c2", "c3", "c4"],
// };

// ---------- Page shell (routes: /chat and /chat/:id) ----------
export default function Chat() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { convs, msgs } = useChat();

  useEffect(() => {
    // if no convo selected, open the most recent one by default
    if (!selectedId && convs.order.length) setSelectedId(convs.order[0]);
  }, [selectedId]);
  // eslint-disable-next-line react-hooks/exhaustive-deps

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
