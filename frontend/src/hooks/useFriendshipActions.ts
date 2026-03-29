import { useChat } from "../context/ChatContext";
import { api } from "../api/api";
import { toConversation, type ConvDto } from "../wire/conversation";
import { dtoToMessage, type MsgDto } from "../wire/message";
import type { CustomResponse } from "../wire/resp";

type CreateConvResponse = {
  conversation: ConvDto;
  message: MsgDto;
};

export function useFriendshipActions() {
  const { friendships, setFriendships, removeConv, loadConvs, loadMsgs } = useChat();

  const acceptFriendship = async (uid: string) => {
    try {
      const resp = await api<CreateConvResponse>(`/friendship/accept`, {
        method: "POST",
        body: JSON.stringify({ uid }),
      });
      const newConv = toConversation(resp.conversation);
      loadConvs([newConv]);
      loadMsgs([dtoToMessage(resp.message)]);
      setFriendships(
        friendships.map((f) =>
          f.uid === uid
            ? { ...f, directConversationId: newConv.id, status: "accepted" }
            : f
        )
      );
    } catch (err) {
      console.error("Error accepting friendship:", err);
    }
  };

  const declineFriendship = async (uid: string) => {
    try {
      await api<CustomResponse>(`/friendship/decline/${uid}`, { method: "DELETE" });
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error declining friendship:", err);
    }
  };

  const cancelRequest = async (uid: string) => {
    try {
      await api<CustomResponse>(`/friendship/cancel/${uid}`, { method: "DELETE" });
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error canceling request:", err);
    }
  };

  const unfriend = async (uid: string, convId: string | null) => {
    try {
      await api<CustomResponse>(`/friendship/delete/${uid}`, { method: "DELETE" });
      if (convId) removeConv(convId);
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error unfriending user:", err);
    }
  };

  const blockUser = async (uid: string, convId: string | null) => {
    try {
      await api<CustomResponse>(`/friendship/block/${uid}`, { method: "POST" });
      if (convId) removeConv(convId);
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error blocking user:", err);
    }
  };

  const addFriend = async (email: string) => {
    try {
      await api<CustomResponse>(`/friendship/request`, {
        method: "POST",
        body: JSON.stringify({ email }),
      });
    } catch (err) {
      console.error("Error sending friend request:", err);
    }
  };

  return { acceptFriendship, declineFriendship, cancelRequest, unfriend, blockUser, addFriend };
}
