package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/utils"
	"github.com/leeryan2000/flashat/wire"
)

type UserHandler struct {
	Repo         repo.UserRepo
	RegisterCode string
	S3Client     *s3.Client
	S3Presigner  *s3.PresignClient
	S3Bucket     string
	S3Region     string
}

type createUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	// ***** replace with email verification code later
	Code string `json:"code"`
}

func (h UserHandler) CreateUser(c *gin.Context) {
	// c.ShouldBindJSON can only read once for a struct body
	var createUserInput createUserInput

	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to bind create user input"})
		return
	}

	// ***** replace with email verification check later
	if createUserInput.Code != h.RegisterCode {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid verification code"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(createUserInput.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		UID:            uuid.New(),
		Name:           createUserInput.Name,
		Email:          createUserInput.Email,
		HashedPassword: hashedPassword,
	}

	if err := h.Repo.CreateUser(c.Request.Context(), &user); err != nil {
		slog.Error("failed to create user", "email", user.Email, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to server"})
		return
	}
	slog.Info("user created", "uid", user.UID, "email", user.Email)
	c.JSON(http.StatusOK, user)
}

func (h UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Repo.GetAllUsers(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h UserHandler) GetUserById(c *gin.Context) {
	uidStr := c.Param("uid")

	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.Repo.GetUserByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) UpdateName(c *gin.Context) {
	uidStr := c.GetString("uid")
	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if err := h.Repo.UpdateName(c.Request.Context(), uid, input.Name); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.GetUserByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) GetAvatarUploadURL(c *gin.Context) {
	uidStr := c.GetString("uid")
	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	key := fmt.Sprintf("avatars/%s.jpg", uid.String())

	presignResult, err := h.S3Presigner.PresignPutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.S3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/jpeg"),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		slog.Error("failed to presign avatar upload url", "uid", uid, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload URL"})
		return
	}

	objectURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s?v=%d", h.S3Bucket, h.S3Region, key, time.Now().Unix())

	c.JSON(http.StatusOK, wire.AvatarUploadURLResponse{
		UploadURL: presignResult.URL,
		ObjectURL: objectURL,
	})
}

func (h UserHandler) SetAvatarURL(c *gin.Context) {
	uidStr := c.GetString("uid")
	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input wire.SetAvatarURLInput
	if err := c.ShouldBindJSON(&input); err != nil || input.AvatarURL == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "avatar_url is required"})
		return
	}

	if err := h.Repo.UpdateAvatarURL(c.Request.Context(), uid, input.AvatarURL); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.GetUserByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) RemoveAvatar(c *gin.Context) {
	uidStr := c.GetString("uid")
	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	key := fmt.Sprintf("avatars/%s.jpg", uid.String())
	if _, err := h.S3Client.DeleteObject(c.Request.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(h.S3Bucket),
		Key:    aws.String(key),
	}); err != nil {
		slog.Error("failed to delete avatar object", "uid", uid, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove photo"})
		return
	}

	if err := h.Repo.UpdateAvatarURL(c.Request.Context(), uid, ""); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.GetUserByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) SendVerification(c *gin.Context) {
	// otp := utils.SendVerification()
	// ***** store otp in redis with expiration associated with user email
	c.JSON(http.StatusOK, gin.H{"message": "Verification sent"})
}
