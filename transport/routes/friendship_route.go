package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
	"github.com/leeryan2000/flashat/server"
)

func InitializeFriendshipRoutes(router *gin.RouterGroup, h *handlers.Handlers, s *server.Server) {
	frh := h.Friendship
	friendshipRoutes := router.Group("/friendship")
	friendshipRoutes.Use(middleware.Authenticate(s))

	friendshipRoutes.POST("/request", frh.RequestFriendship)
	friendshipRoutes.POST("/accept", frh.AcceptFriendship)

	friendshipRoutes.DELETE("/delete/:friend_uid", frh.DeleteFriendship)
	friendshipRoutes.DELETE("/decline/:friend_uid", frh.DeclineFriendship)
	friendshipRoutes.DELETE("/cancel/:friend_uid", frh.CancelFriendshipRequest)

	friendshipRoutes.POST("/block/:uid", frh.BlockUser)

	friendshipRoutes.GET("", frh.ListFriendships)
}
