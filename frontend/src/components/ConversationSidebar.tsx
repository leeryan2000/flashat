import { useMemo, useState } from "react";
import type { ConvState } from "../context/ChatContext";
import type { Conversation } from "../wire/conversation";


function formatConvTime(ts: number) {
  const now = new Date();
  const date = new Date(ts);
  const diffMins = Math.max(1, Math.round((Date.now() - ts) / 60000));

  if (diffMins < 60) return `${diffMins}m`;
  if (diffMins < 1440) return `${Math.round(diffMins / 60)}h`;

  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  }
  return date.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
}

export default function ConversationsSidebar({ convs, activeConvId, onOpen }:{ convs: ConvState; activeConvId?: string; onOpen:(id:string)=>void }) {
  // Query for the title search
  const [query, setQ] = useState("");
  
  // filters the conversations by title
  const filtered = useMemo(() => {
    const needle = query.trim().toLowerCase();
    return convs.order
    .map(id => convs.entities[id])
    .filter((c): c is Conversation => !!c) 
    .filter(conv => conv.title?.toLowerCase().includes(needle));

  }, [convs, query]);

  return (
    <div className="h-full text-slate-100 w-80 flex-shrink-0 grid grid-rows-[auto_1fr]" style={{ background: "var(--sidebar-bg)" }}>
      {/* search */}
      <div className="p-3 border-b" style={{ borderColor: "var(--sidebar-item)" }}>
        <input
          value={query}
          onChange={e=>setQ(e.target.value)}
          placeholder="Search conversations"
          className="w-full rounded-xl px-3 py-2 text-sm text-slate-100 placeholder:text-slate-400 outline-none focus:ring-2"
          style={{ background: "var(--sidebar-item)", outlineColor: "var(--sidebar-ring)" }}
        />
      </div>

      {/* list */}
      <div className="overflow-y-auto">
        {filtered.map(conv => (
          <button
            key={conv.id}
            onClick={()=>onOpen(conv.id)}
            className="w-full text-left px-3 py-3 gap-3 flex items-start transition-colors duration-100"
            style={{ background: activeConvId===conv.id ? "var(--sidebar-item)" : undefined }}
            onMouseEnter={e => { if (activeConvId !== conv.id) (e.currentTarget as HTMLButtonElement).style.background = "var(--sidebar-item)"; }}
            onMouseLeave={e => { if (activeConvId !== conv.id) (e.currentTarget as HTMLButtonElement).style.background = ""; }}
          >
            <div className="h-10 w-10 rounded-full flex items-center justify-center text-white text-xs shrink-0" style={{ background: "var(--primary)" }}>{conv.title.slice(0,2)}</div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate font-medium">
                  {conv.title}{conv.type === "group" && <span className="text-slate-400 font-normal"> ({conv.participants.length})</span>}
                </p>
                <span className="text-xs text-slate-400 ml-auto">
                  {typeof conv.lastMsgTs === "number" ? formatConvTime(conv.lastMsgTs) : "--"}
                </span>
              </div>
              <p className="truncate text-sm text-slate-300">{conv.lastMsgText}</p>
            </div>
            {!!conv.unreadCount && (
              <span className="ml-2 inline-flex items-center justify-center rounded-full bg-emerald-500 text-white text-xs px-2 py-0.5">
                {conv.unreadCount}
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