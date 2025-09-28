import { useState } from "react";

export function Composer({ convId }: { convId: string }) {
  const [text, setText] = useState("");

  function send() {
    if (!text.trim()) return;
    // TODO: replace with your ws/http send
    alert(`Send to ${convId}:\n${text}`);
    setText("");
  }

  return (
    <div className="mt-3 grid grid-cols-[1fr_auto] gap-2 rounded-xl border border-slate-200 p-2 bg-white">
      <textarea
        rows={1}
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            send();
          }
        }}
        placeholder="Type a message… (Enter to send, Shift+Enter for new line)"
        className="resize-none w-full px-3 py-2 rounded-lg outline-none focus:ring-2 focus:ring-indigo-400"
      />
      <button
        onClick={send}
        className="px-4 py-2 rounded-lg bg-indigo-600 text-white font-medium hover:brightness-110 active:scale-95 transition"
      >
        Send
      </button>
    </div>
  );
}