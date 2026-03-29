import { useChat } from "../context/ChatContext";
import { api } from "../api/api";
import type { CustomResponse } from "../wire/resp";

export function useFriendshipActions() {
  const { friendships, setFriendships, removeConv } = useChat();

  const blockUser = async (uid: string, convId: string | null) => {
    try {
      await api<CustomResponse>(`/friendship/block/${uid}`, { method: "POST" });
      if (convId) removeConv(convId);
      setFriendships(
        friendships.map((f) =>
          f.uid === uid
            ? { ...f, status: "blocked", direction: "outgoing", directConversationId: null }
            : f
        )
      );
    } catch (err) {
      console.error("Error blocking user:", err);
    }
  };

  return { blockUser };
}
