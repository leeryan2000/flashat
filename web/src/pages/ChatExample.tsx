import React, { useMemo, useRef, useState, useEffect } from "react";
import { Link, Outlet, useNavigate, useParams } from "react-router-dom";

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
export interface Conversation {
  id: string;
  title: string;      // display name / group name
  last: string;       // last message preview
  ts: number;         // unix ms
  unread?: number;    // unread count
}

export interface Message {
  id: string;
  convoId: string;
  from: string;       // uid or display name
  text: string;
  ts: number;         // unix ms
  mine?: boolean;     // render bubble on the right
}

// ---------- Mock data (replace with your API) ----------
const MOCK_CONVOS: Conversation[] = [
  { id: "c1", title: "Alice", last: "See you tomorrow!", ts: Date.now() - 1000 * 60 * 3, unread: 2 },
  { id: "c2", title: "Bob • Project", last: "Pushed a new branch.", ts: Date.now() - 1000 * 60 * 11 },
  { id: "c3", title: "Team Rocket", last: "Prepare for trouble.", ts: Date.now() - 1000 * 60 * 60 * 5, unread: 7 },
  { id: "c4", title: "Huai'en (You)", last: "Drafting the layout now", ts: Date.now() - 1000 * 60 * 60 * 10 },
];

const MOCK_MESSAGES: Message[] = [
  { id: "m1", convoId: "c1", from: "Alice", text: "Hey!", ts: Date.now() - 1000 * 60 * 60, mine: false },
  { id: "m2", convoId: "c1", from: "Huai'en", text: "Yo!", ts: Date.now() - 1000 * 60 * 58, mine: true },
  { id: "m3", convoId: "c1", from: "Alice", text: "See you tomorrow!", ts: Date.now() - 1000 * 60 * 3, mine: false },
  { id: "m4", convoId: "c2", from: "Bob", text: "Pushed a new branch.", ts: Date.now() - 1000 * 60 * 11, mine: false },
  { id: "m5", convoId: "c3", from: "Jessie", text: "Prepare for trouble.", ts: Date.now() - 1000 * 60 * 60 * 5, mine: false },
  { id: "m6", convoId: "c3", from: "James", text: "And make it double!", ts: Date.now() - 1000 * 60 * 60 * 5 + 10000, mine: false },
  { id: "m7", convoId: "c4", from: "Huai'en", text: "Drafting the layout now", ts: Date.now() - 1000 * 60 * 60 * 10, mine: true },
];

// ---------- Helpers ----------
function timeAgo(ts: number) {
  const mins = Math.max(1, Math.round((Date.now() - ts) / 60000));
  if (mins < 60) return `${mins}m`;
  const hrs = Math.round(mins / 60);
  if (hrs < 24) return `${hrs}h`;
  return `${Math.round(hrs / 24)}d`;
}

// ---------- Sidebar (search + list) ----------
function ConversationsSidebar({ items, activeId, onOpen }:{ items: Conversation[]; activeId?: string; onOpen:(id:string)=>void }) {
  const [q, setQ] = useState("");
  const filtered = useMemo(
    () => items.filter(c => c.title.toLowerCase().includes(q.trim().toLowerCase())),
    [items, q]
  );

  return (
    <div className="h-full bg-slate-900 text-slate-100 w-80 flex-shrink-0 grid grid-rows-[auto_1fr]">
      {/* search */}
      <div className="p-3 border-b border-slate-800">
        <input
          value={q}
          onChange={e=>setQ(e.target.value)}
          placeholder="Search conversations"
          className="w-full rounded-xl bg-slate-800/60 px-3 py-2 text-sm placeholder:text-slate-400 outline-none focus:ring-2 focus:ring-indigo-400"
        />
      </div>

      {/* list */}
      <div className="overflow-y-auto">
        {filtered.map(c => (
          <button
            key={c.id}
            onClick={()=>onOpen(c.id)}
            className={[
              "w-full text-left px-3 py-3 gap-3 flex items-start",
              activeId===c.id ? "bg-slate-800" : "hover:bg-slate-800/60"
            ].join(" ")}
          >
            <div className="h-10 w-10 rounded-full bg-slate-700 flex items-center justify-center text-xs">{c.title.slice(0,2)}</div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate font-medium">{c.title}</p>
                <span className="text-xs text-slate-400 ml-auto">{timeAgo(c.ts)}</span>
              </div>
              <p className="truncate text-sm text-slate-300">{c.last}</p>
            </div>
            {!!c.unread && (
              <span className="ml-2 inline-flex items-center justify-center rounded-full bg-emerald-500 text-white text-xs px-2 py-0.5">
                {c.unread}
              </span>
            )}
          </button>
        ))}
        {filtered.length===0 && (
          <div className="p-6 text-sm text-slate-400">No conversations</div>
        )}
      </div>
    </div>
  );
}

// ---------- Messages pane ----------
function MessagesPane({ convoId }:{ convoId: string }) {
  const list = useMemo(() => MOCK_MESSAGES.filter(m => m.convoId===convoId).sort((a,b)=>a.ts-b.ts), [convoId]);
  const boxRef = useRef<HTMLDivElement|null>(null);

  useEffect(()=>{
    // auto scroll to bottom on convo change or new message
    boxRef.current?.scrollTo({ top: boxRef.current.scrollHeight })
  }, [convoId, list.length]);

  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      {/* messages list */}
      <div ref={boxRef} className="overflow-y-auto space-y-2 pr-2">
        {list.map(m => (
          <div key={m.id} className={["flex", m.mine?"justify-end":"justify-start"].join(" ") }>
            <div className={[
              "max-w-[75%] rounded-2xl px-3 py-2 text-sm shadow",
              m.mine ? "bg-indigo-600 text-white rounded-br-sm" : "bg-slate-100 text-slate-900 rounded-bl-sm"
            ].join(" ")}>
              <p className="whitespace-pre-wrap break-words">{m.text}</p>
              <p className="mt-1 text-[11px] opacity-70">{new Date(m.ts).toLocaleTimeString()}</p>
            </div>
          </div>
        ))}
        {list.length===0 && (
          <div className="h-full grid place-items-center text-slate-400">No messages yet.</div>
        )}
      </div>

      {/* composer */}
      <Composer convoId={convoId} />
    </div>
  );
}

function Composer({ convoId }:{ convoId: string }) {
  const [text, setText] = useState("");

  function send() {
    if (!text.trim()) return;
    // TODO: replace with your ws/http send
    alert(`Send to ${convoId}:\n${text}`);
    setText("");
  }

  return (
    <div className="mt-3 grid grid-cols-[1fr_auto] gap-2 rounded-xl border border-slate-200 p-2 bg-white">
      <textarea
        rows={1}
        value={text}
        onChange={e=>setText(e.target.value)}
        onKeyDown={(e)=>{
          if (e.key==="Enter" && !e.shiftKey){ e.preventDefault(); send(); }
        }}
        placeholder="Type a message… (Enter to send, Shift+Enter for new line)"
        className="resize-none w-full px-3 py-2 rounded-lg outline-none focus:ring-2 focus:ring-indigo-400"
      />
      <button onClick={send} className="px-4 py-2 rounded-lg bg-indigo-600 text-white font-medium hover:brightness-110 active:scale-95 transition">
        Send
      </button>
    </div>
  );
}

// ---------- Page shell (routes: /chat and /chat/:id) ----------
export default function Chat() {
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const convos = useMemo(
    () => [...MOCK_CONVOS].sort((a,b)=> b.ts - a.ts),
    []
  );

  useEffect(()=>{
    // if no convo selected, open the most recent one by default
    if (!selectedId && convos.length) setSelectedId(convos[0].id);
  }, [selectedId]);
    // eslint-disable-next-line react-hooks/exhaustive-deps

  return (
    <div className="grid md:grid-cols-[320px_1fr] min-h-screen">
      <ConversationsSidebar
        items={MOCK_CONVOS}
        activeId={selectedId ?? undefined}
        onOpen={(cid) => setSelectedId(cid)}   // <-- no navigate
      />
      <section className="bg-white p-6">
        {selectedId ? (
          <MessagesPane convoId={selectedId} />
        ) : (
          <div className="h-full grid place-items-center text-slate-400">
            Pick a conversation on the left
          </div>
        )}
      </section>
    </div>
  );
}
