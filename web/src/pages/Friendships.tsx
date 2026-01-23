import { useChat } from "../context/ChatContext";
import FriendshipList from "../components/FriendshipList";
import { useNavigate } from "react-router-dom";
import { PATHS } from "../routes/paths";
import { api } from "../api/api";
import {
  toConversation as dtoToConversation,
  type ConvDto,
} from "../wire/conversation";
import { dtoToMessage, type MsgDto } from "../wire/message";
import type { Friendship } from "../wire/friendship";
import type { CustomResponse } from "../wire/resp";

type AcceptFriendship = {
  conversation: ConvDto;
  message: MsgDto;
};

export default function Friendships() {
  const {
    friendships,
    setFriendships,
    setActiveConvId,
    loadConvs,
    loadMsgs,
    removeConv,
  } = useChat();
  const navigate = useNavigate();

  const handleChat = (convId: string) => {
    console.log("Navigating to chat with convId:", convId);
    setActiveConvId(convId);
    navigate(PATHS.chat);
  };

  const handleAccept = async (uid: string) => {
    try {
      const resp = await api<AcceptFriendship>(`/friendship/accept`, {
        method: "POST",
        body: JSON.stringify({
          uid,
        }),
      });

      const newDirectConv = dtoToConversation(resp.conversation);
      loadConvs([newDirectConv]);
      const welcomeMsg = dtoToMessage(resp.message);
      loadMsgs([welcomeMsg]);

      // update status for friendship list
      const updatedList: Friendship[] = friendships.map((friend) => {
        if (friend.uid === uid) {
          return {
            ...friend,
            directConversationId: newDirectConv.id,
            status: "accepted",
          };
        }
        return friend;
      });
      setFriendships(updatedList);
    } catch (err: any) {
      console.error("Error accepting friendship:", err);
    }
  };

  const onDecline = async (uid: string) => {
    try {
        const resp = await api<CustomResponse>(`/friendship/decline/${uid}`, {
          method: "DELETE" 
        });
      
        const updatedList: Friendship[] = friendships.filter(
          (friend) => friend.uid !== uid
        );
        setFriendships(updatedList)

        console.log(resp.message)
    } catch (err: any) {
      console.error("Error declining friendship:", err);
    }
  };

  const onCancel = async (uid: string) => {};

  const onUnfriend = async (uid: string, convId: string | null) => {
    try {
      const resp = await api<CustomResponse>(`/friendship/delete/${uid}`, {
        method: "DELETE",
      });

      if (convId) {
        removeConv(convId);
      }
      const updatedList: Friendship[] = friendships.filter(
        (friend) => friend.uid !== uid
      );
      setFriendships(updatedList);

      console.log(resp.message)
    } catch (err: any) {
      console.error("Error unfriending user:", err);
    }
  };

  return (
    <div className="flex h-screen bg-black">
      <main className="flex-1 h-full">
        <FriendshipList
          data={friendships}
          onChatClick={handleChat}
          onAccept={handleAccept}
          onDecline={onDecline}
          onCancel={onCancel}
          onUnfriend={onUnfriend}
        />
      </main>
    </div>
  );
}
