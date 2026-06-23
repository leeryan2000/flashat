import { X, Shield, User } from "lucide-react";
import type { Participant } from "../wire/conversation";

type ParticipantsPanelProps = {
  participants: Participant[];
  onClose: () => void;
};

export function ParticipantsPanel({ participants, onClose }: ParticipantsPanelProps) {
  const admins = participants.filter(p => p.role === "admin" || p.role === "creator");
  const members = participants.filter(p => p.role === "member");

  const renderRow = (p: Participant) => (
    <div key={p.uid} className="flex items-center gap-3 px-4 py-2.5 transition" onMouseEnter={e => (e.currentTarget.style.background = "color-mix(in srgb, var(--primary) 8%, transparent)")} onMouseLeave={e => (e.currentTarget.style.background = "")}>
      <div className="h-8 w-8 rounded-full flex items-center justify-center text-white text-sm font-semibold shrink-0" style={{ background: "var(--primary)" }}>
        {p.name.slice(0, 2).toUpperCase()}
      </div>
      <span className="flex-1 text-sm text-slate-200 truncate">{p.name}</span>
      {p.role === "creator" ? (
        <span className="flex items-center gap-1 text-[11px] font-medium text-amber-400 bg-amber-400/10 px-2 py-0.5 rounded-full">
          <Shield size={11} /> Creator
        </span>
      ) : p.role === "admin" ? (
        <span className="flex items-center gap-1 text-[11px] font-medium px-2 py-0.5 rounded-full" style={{ color: "var(--primary)", background: "color-mix(in srgb, var(--primary) 12%, transparent)" }}>
          <Shield size={11} /> Admin
        </span>
      ) : (
        <span className="flex items-center gap-1 text-[11px] font-medium text-slate-500 bg-slate-700/50 px-2 py-0.5 rounded-full">
          <User size={11} /> Member
        </span>
      )}
    </div>
  );

  return (
    <div className="absolute inset-y-0 right-0 w-64 border-l flex flex-col z-30 shadow-xl" style={{ background: "var(--sidebar-item)", borderColor: "color-mix(in srgb, var(--primary) 15%, transparent)" }}>
      <div className="flex items-center justify-between px-4 py-3 border-b" style={{ borderColor: "color-mix(in srgb, var(--primary) 15%, transparent)" }}>
        <span className="text-sm font-semibold text-slate-200">Members · {participants.length}</span>
        <button onClick={onClose} className="text-slate-400 hover:text-white transition"><X size={18} /></button>
      </div>
      <div className="flex-1 overflow-y-auto py-2">
        {admins.length > 0 && (
          <div>
            <p className="px-4 pt-2 pb-1 text-[11px] font-semibold uppercase tracking-wider text-slate-500">Admins — {admins.length}</p>
            {admins.map(renderRow)}
          </div>
        )}
        {members.length > 0 && (
          <div>
            <p className="px-4 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wider text-slate-500">Members — {members.length}</p>
            {members.map(renderRow)}
          </div>
        )}
      </div>
    </div>
  );
}
