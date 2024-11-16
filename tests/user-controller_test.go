package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/dewciu/f1_api/pkg/models"
	"github.com/dewciu/f1_api/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	token  string
}

func (suite *UserControllerTestSuite) SetupSuite() {
	// Setup the test environment
	tablesAffected := []string{"users"}
	setupDB(tablesAffected)
	suite.router = routes.SetupRouter()
	suite.token = suite.authenticate()
	fmt.Println("SetupSuite")
}

func (suite *UserControllerTestSuite) authenticate() string {
	// Create a new user for authentication
	user := m.User{
		Username: "authuser",
		Email:    "authuser@email.com",
		Password: "authpassword",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Perform login to get the token
	loginData := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}
	loginJson, _ := json.Marshal(loginData)
	loginReq, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	suite.router.ServeHTTP(loginW, loginReq)

	var response map[string]string
	json.Unmarshal(loginW.Body.Bytes(), &response)

	return response["token"]
}

func (suite *UserControllerTestSuite) TestCreateUser(t *testing.T) {
	// Create a new user
	user := m.User{
		Username: "testuser",
		Email:    "testemail@email.com",
		Password: "testpassword",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser m.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.Nil(t, err)
	assert.Equal(t, user.Username, createdUser.Username)
	assert.Equal(t, user.Email, createdUser.Email)
	fmt.Println("TestCreateUser")
}

func (suite *UserControllerTestSuite) TestGetUser(t *testing.T) {
	// Get an existing user
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user m.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.Equal(t, "authuser", user.Username)
	assert.Equal(t, "authuser@email.com", user.Email)
	fmt.Println("TestGetUser")
}

func (suite *UserControllerTestSuite) TestUpdateUser(t *testing.T) {
	// Update an existing user
	updatedData := map[string]string{
		"username": "updateduser",
		"email":    "updatedemail@email.com",
	}
	jsonData, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedUser m.User
	err := json.Unmarshal(w.Body.Bytes(), &updatedUser)
	assert.Nil(t, err)
	assert.Equal(t, updatedData["username"], updatedUser.Username)
	assert.Equal(t, updatedData["email"], updatedUser.Email)
	fmt.Println("TestUpdateUser")
}

func (suite *UserControllerTestSuite) TestDeleteUser(t *testing.T) {
	// Delete an existing user
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	fmt.Println("TestDeleteUser")
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
