package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/mocks"
	"github.com/leeryan2000/flashat/models"
	"github.com/stretchr/testify/assert"
)

func setupRouter(h handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/createUser", h.CreateUser)
	r.GET("/getAllUsers", h.GetAllUsers)
	r.GET("/getUserById/:id", h.GetUserById)
	return r
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepoMock)
	user := models.User{Email: "testEmail@email.com", Hashed_Password: "123"}
	mockRepo.On("CreateUser", &user).Return(nil)

	h := handlers.UserHandler{Repo: mockRepo}
	r := setupRouter(h)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/createUser", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepoMock)
	users := []models.User{{Email: "user1@example.com"}, {Email: "user2@example.com"}}
	mockRepo.On("GetAllUsers").Return(users, nil)

	h := handlers.UserHandler{Repo: mockRepo}
	r := setupRouter(h)

	req, _ := http.NewRequest("GET", "/getAllUsers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "user1@example.com", response[0].Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepoMock)
	user := &models.User{Id: 1, Email: "alice@example.com"}
	mockRepo.On("GetUserById", uint(1)).Return(user, nil)

	h := handlers.UserHandler{Repo: mockRepo}
	r := setupRouter(h)

	req, _ := http.NewRequest("GET", "/getUserById/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
