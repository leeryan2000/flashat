import { useChat } from "../context/ChatContext";
import FriendsList from "../components/FriendsList";

export default function Friendships() {
    const { friendships } = useChat();
    const handleChat = (convId: string) => {
        
    }

    const handleAccept = async (uid: string) => {
    }

    const onDecline = async (uid: string) => {
    }

    const onCancel = async (uid: string) => {
    }
    return <div className="flex h-screen bg-black">
      <main className="flex-1 h-full">
        <FriendsList 
           data={friendships}
           onChatClick={handleChat}
           onAccept={handleAccept}
           onDecline={(uid) => console.log("Declined", uid)}
           onCancel={(uid) => console.log("Cancelled", uid)}
        />
      </main>
    </div>
}