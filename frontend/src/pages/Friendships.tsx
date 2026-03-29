import { useChat } from "../context/ChatContext";
import FriendshipList from "../components/FriendshipList";
import { useNavigate } from "react-router-dom";
import { PATHS } from "../routes/paths";
import { api } from "../api/api";
import { toConversation, type CreateConvResponse } from "../wire/conversation";
import { dtoToMessage } from "../wire/message";
import { useFriendshipActions } from "../hooks/useFriendshipActions";

export default function Friendships() {
  const { friendships, setActiveConvId, loadConvs, loadMsgs } = useChat();
  const { acceptFriendship, declineFriendship, cancelRequest, unfriend, blockUser, unblockUser, addFriend } =
    useFriendshipActions();
  const navigate = useNavigate();

  const onChatClick = (convId: string) => {
    setActiveConvId(convId);
    navigate(PATHS.chat);
  };

  const onCreateGroup = async (name: string, participants: string[]) => {
    try {
      const resp = await api<CreateConvResponse>(`/conversation/group`, {
        method: "POST",
        body: JSON.stringify({ name, participants }),
      });
      loadConvs([toConversation(resp.conversation)]);
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
          onAccept={acceptFriendship}
          onDecline={declineFriendship}
          onCancel={cancelRequest}
          onUnfriend={unfriend}
          onBlock={blockUser}
          onUnblock={unblockUser}
          onAddFriend={addFriend}
          onCreateGroup={onCreateGroup}
        />
      </main>
    </div>
  );
}
