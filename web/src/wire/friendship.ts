export type Friendship = {
    uid: string,
    name: string,
    email: string,
    avatarUrl: string,
    directConversationId: string | null,  // null if the friendship hasn't been accepted yet
    status: "pending" | "accepted" | "blocked",
    direction: "incoming" | "outgoing",
}

export type FriendshipDto = {
    uid: string,
    name: string,
    email: string,
    avatar_url: string,
    direct_conversation_id: string | null,
    status: "pending" | "accepted" | "blocked",
    direction: "incoming" | "outgoing",
}


export const toFriendship = (w: FriendshipDto): Friendship => ({
    uid: w.uid,
    name: w.name,
    email: w.email,
    avatarUrl: w.avatar_url,
    directConversationId: w.direct_conversation_id,
    status: w.status,
    direction: w.direction,
}); 

