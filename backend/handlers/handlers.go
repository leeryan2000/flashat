package handlers

type Handlers struct {
	User         UserHandler
	Auth         AuthHandler
	Websocket    WebsocketHandler
	Conversation ConversationHandler
	Message      MessageHandler
	Friendship   FriendshipHandler
}
