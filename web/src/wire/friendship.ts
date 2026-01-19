export type Friendship = {
    uid: string,
    name: string,
    email: string,
    avatarUrl?: string | null,
    directConversationId?: string,
    status: "pending" | "accepted" | "blocked",
    direction: "incoming" | "outgoing",
}

export type FriendshipDto = {
    uid: string,
    name: string,
    email: string,
    avatar_url?: string | null,
    direct_conversation_id?: string,
    status: "pending" | "accepted" | "blocked",
    direction: "incoming" | "outgoing",
}


export const toFriendship = (w: FriendshipDto): Friendship => ({
    uid: w.uid,
    name: w.name,
    email: w.email,
    avatarUrl: w.avatar_url ?? null,
    directConversationId: w.direct_conversation_id,
    status: w.status,
    direction: w.direction,
}); 

