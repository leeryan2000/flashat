package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
	"github.com/leeryan2000/flashat-posts/models"
	"github.com/leeryan2000/flashat-posts/repo"
)

type PostHandler struct {
	Repo repo.PostRepo
	// AuthClient is used to look up the requester's friends when building
	// the feed — the posts service has no direct access to the
	// friendships table (separate database, separate service).
	AuthClient authpb.AuthInternalClient
}

type createPostRequest struct {
	Body string `json:"body" binding:"required"`
}

type addCommentRequest struct {
	Body string `json:"body" binding:"required"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString("uid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "body is required"})
		return
	}

	post, err := h.Repo.CreatePost(c.Request.Context(), uid, req.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, post)
}

// ListFeed scopes the feed to the requester + their accepted friends —
// no edit/delete, no third-party author info here (see the frontend
// plan: display info is resolved client-side from data it already has).
func (h *PostHandler) ListFeed(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString("uid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	friendResp, err := h.AuthClient.GetFriendUIDs(c.Request.Context(), &authpb.GetFriendUIDsRequest{Uid: uid.String()})
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to resolve friends"})
		return
	}

	visibleUIDs := make([]uuid.UUID, 0, len(friendResp.GetFriendUids())+1)
	visibleUIDs = append(visibleUIDs, uid)
	for _, s := range friendResp.GetFriendUids() {
		friendUID, err := uuid.Parse(s)
		if err != nil {
			continue
		}
		visibleUIDs = append(visibleUIDs, friendUID)
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 50
	}

	var posts []models.Post
	if beforeSeqStr := c.Query("before_seq"); beforeSeqStr != "" {
		beforeSeq, err := strconv.ParseInt(beforeSeqStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid before_seq"})
			return
		}
		posts, err = h.Repo.ListFeedBefore(c.Request.Context(), uid, visibleUIDs, beforeSeq, limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		posts, err = h.Repo.ListFeedLatest(c.Request.Context(), uid, visibleUIDs, limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString("uid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.Repo.GetPost(c.Request.Context(), uid, postID)
	if err != nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}

	c.JSON(200, post)
}

func (h *PostHandler) ToggleLike(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString("uid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	liked, likesCount, err := h.Repo.ToggleLike(c.Request.Context(), postID, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"liked": liked, "likes_count": likesCount})
}

func (h *PostHandler) AddComment(c *gin.Context) {
	uid, err := uuid.Parse(c.GetString("uid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	var req addCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "body is required"})
		return
	}

	comment, err := h.Repo.AddComment(c.Request.Context(), postID, uid, req.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, comment)
}

func (h *PostHandler) ListComments(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 50
	}

	comments, err := h.Repo.ListComments(c.Request.Context(), postID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, comments)
}
