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
	friendshipRoutes.DELETE("/reject/:friend_uid", frh.RejectFriendship)

	friendshipRoutes.GET("/", frh.ListFriendships)
	friendshipRoutes.GET("/requests", frh.ListFriendshipRequests)
}
