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
    <div className="p-3 bg-white">
      <div className="flex items-end gap-2 rounded-2xl border border-slate-200 shadow-sm bg-white px-3 py-2">
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
          className="resize-none flex-1 text-base text-slate-800 placeholder:text-slate-400 outline-none overflow-y-auto py-2 leading-6"
          style={{ height: "40px", maxHeight: "160px" }}
        />
        <button
          onClick={send}
          disabled={text.trim() === ""}
          className="flex items-center gap-2 px-4 py-2 rounded-xl bg-indigo-600 text-white text-sm font-medium hover:bg-indigo-500 active:scale-95 transition disabled:opacity-40 disabled:cursor-not-allowed shrink-0"
        >
          <Send size={15} />
          Send
        </button>
      </div>
    </div>
  );
}
