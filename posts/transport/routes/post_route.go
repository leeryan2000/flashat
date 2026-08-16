package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat-posts/handlers"
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
	"github.com/leeryan2000/flashat-posts/middleware"
)

// InitializePostRoutes registers routes on router, which is already
// scoped to /api/posts (see routes.go) — nginx only ever forwards
// /api/posts/* to this service, so there's no separate /posts subgroup
// the way backend nests /api/message under /api.
func InitializePostRoutes(router *gin.RouterGroup, h *handlers.Handlers, authClient authpb.AuthInternalClient) {
	ph := h.Post

	authed := router.Group("")
	authed.Use(middleware.Authenticate(authClient))

	authed.POST("", ph.CreatePost)
	authed.GET("/feed", ph.ListFeed)
	authed.GET("/:post_id", ph.GetPost)
	authed.POST("/:post_id/like", ph.ToggleLike)
	authed.POST("/:post_id/comments", ph.AddComment)
	authed.GET("/:post_id/comments", ph.ListComments)
}
