export type Friendship = {
    uid: string,
    name: string,
    email: string,
    avatarUrl?: string | null,
    directConversationId?: string,
}

export type FriendshipDto = {
    uid: string,
    name: string,
    email: string,
    avatar_url?: string | null,
    direct_conversation_id?: string,
}


export const toFriendship = (w: FriendshipDto): Friendship => ({
    uid: w.uid,
    name: w.name,
    email: w.email,
    avatarUrl: w.avatar_url ?? null,
    directConversationId: w.direct_conversation_id,
}); 

