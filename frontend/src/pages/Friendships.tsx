import { useChat } from "../context/ChatContext";
import FriendshipList from "../components/FriendshipList";
import { useNavigate } from "react-router-dom";
import { PATHS } from "../routes/paths";
import { api } from "../api/api";
import { toConversation, type CreateConvResponse } from "../wire/conversation";
import { dtoToMessage } from "../wire/message";
import { toFriendship, type Friendship, type FriendshipDto } from "../wire/friendship";
import type { CustomResponse } from "../wire/resp";
import { useFriendshipActions } from "../hooks/useFriendshipActions";
import { useAuth } from "../context/AuthContext";

export default function Friendships() {
  const { friendships, setFriendships, setActiveConvId, loadConvs, loadMsgs, removeConv } = useChat();
  const { blockUser } = useFriendshipActions();
  const { user } = useAuth();
  const navigate = useNavigate();

  const onChatClick = (convId: string) => {
    setActiveConvId(convId);
    navigate(PATHS.chat);
  };

  const onAccept = async (uid: string) => {
    try {
      const resp = await api<CreateConvResponse>(`/friendship/accept`, {
        method: "POST",
        body: JSON.stringify({ uid }),
      });
      const newConv = toConversation(resp.conversation);
      loadConvs([newConv]);
      loadMsgs([dtoToMessage(resp.message)]);
      setFriendships(
        friendships.map((f): Friendship =>
          f.uid === uid ? { ...f, directConversationId: newConv.id, status: "accepted" } : f
        )
      );
    } catch (err) {
      console.error("Error accepting friendship:", err);
    }
  };

  const onDecline = async (uid: string) => {
    try {
      await api<CustomResponse>(`/friendship/decline/${uid}`, { method: "DELETE" });
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error declining friendship:", err);
    }
  };

  const onCancel = async (uid: string) => {
    try {
      await api<CustomResponse>(`/friendship/cancel/${uid}`, { method: "DELETE" });
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error canceling friendship request:", err);
    }
  };

  const onUnfriend = async (uid: string, convId: string | null) => {
    try {
      await api<CustomResponse>(`/friendship/delete/${uid}`, { method: "DELETE" });
      if (convId) removeConv(convId);
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error unfriending user:", err);
    }
  };

  const onUnblock = async (uid: string) => {
    try {
      await api<CustomResponse>(`/friendship/delete/${uid}`, { method: "DELETE" });
      setFriendships(friendships.filter((f) => f.uid !== uid));
    } catch (err) {
      console.error("Error unblocking user:", err);
    }
  };

  const onAddFriend = async (email: string) => {
    try {
      const dto = await api<FriendshipDto>(`/friendship/request`, {
        method: "POST",
        body: JSON.stringify({ email }),
      });
      setFriendships([...friendships, toFriendship(dto)]);
    } catch (err) {
      console.error("Error sending friend request:", err);
      throw err;
    }
  };

  const onCreateGroup = async (name: string, participantIds: string[]) => {
    try {
      const resp = await api<CreateConvResponse>(`/conversation/group`, {
        method: "POST",
        body: JSON.stringify({ name, participantIds }),
      });

      const selectedMembers = friendships
        .filter(f => participantIds.includes(f.uid))
        .map(f => ({ uid: f.uid, name: f.name, role: "member" as const }));

      const conv = {
        ...toConversation(resp.conversation),
        participants: [
          { uid: user!.uid, name: user!.name, role: "creator" as const },
          ...selectedMembers,
        ],
      };

      loadConvs([conv]);
      loadMsgs([dtoToMessage(resp.message)]);
    } catch (err) {
      console.error("Error creating group conversation:", err);
    }
  };

  return (
    <div className="flex h-screen bg-black">
      <main className="flex-1 h-full">
        <FriendshipList
          data={friendships}
          onChatClick={onChatClick}
          onAccept={onAccept}
          onDecline={onDecline}
          onCancel={onCancel}
          onUnfriend={onUnfriend}
          onBlock={blockUser}
          onUnblock={onUnblock}
          onAddFriend={onAddFriend}
          onCreateGroup={onCreateGroup}
        />
      </main>
    </div>
  );
}
