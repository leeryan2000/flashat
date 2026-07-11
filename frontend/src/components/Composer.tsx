import { useState } from "react";
import { useChat } from "../context/ChatContext";
import { Send } from "lucide-react";

export function Composer({ convId }: { convId: string }) {
  const [text, setText] = useState("");
  const { sendMessage } = useChat();

  function send() {
    if (text.trim() === "") return;
    sendMessage(text, convId);
    setText("");
  }

  return (
    <div className="p-3" style={{ background: "var(--chat-bg)" }}>
      <div
        className="flex items-end gap-2 rounded-2xl border px-3 py-1.5"
        style={{ background: "var(--panel)", borderColor: "var(--panel-border)" }}
      >
        <textarea
          rows={1}
          value={text}
          onChange={(e) => {
            setText(e.target.value);
            e.target.style.height = "auto";
            e.target.style.height = `${e.target.scrollHeight}px`;
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              send();
            }
          }}
          placeholder="Write a message..."
          className="resize-none flex-1 text-[15px] bg-transparent outline-none overflow-y-auto py-2 leading-6 placeholder:text-[color:var(--text-faint)]"
          style={{ height: "40px", maxHeight: "160px", color: "var(--text)" }}
        />
        <button
          onClick={send}
          disabled={text.trim() === ""}
          title="Send"
          className="h-9 w-9 mb-0.5 rounded-full flex items-center justify-center text-white active:scale-95 transition disabled:opacity-40 disabled:cursor-not-allowed shrink-0"
          style={{ background: "var(--primary)" }}
        >
          <Send size={16} className="-ml-0.5" />
        </button>
      </div>
    </div>
  );
}
