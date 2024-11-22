package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/dewciu/f1_api/pkg/models"
	"github.com/dewciu/f1_api/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type UserCreateTestSuite struct {
	suite.Suite
	db           *gorm.DB
	pgContainter tc.Container
	ctx          context.Context
	router       *gin.Engine
	baseHeader   http.Header
}

func (suite *UserCreateTestSuite) SetupSuite() {
	// Setup the test environment
	suite.db, suite.pgContainter, suite.ctx = SetupDB([]string{"users"})
	suite.router = routes.SetupRouter(suite.db)
	fmt.Println(suite.db)
	token := suite.authenticate()
	fmt.Println("tokee")
	fmt.Println(token)
	suite.baseHeader = http.Header{
		"Authorization": []string{token},
		"Content-Type":  []string{"application/json"},
	}
}

func (suite *UserCreateTestSuite) authenticate() string {
	loginData := map[string]string{
		"password": "admin",
		"username": "admin",
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

func (suite *UserCreateTestSuite) TestSuccessfulUserCreation() {
	var createdUser m.User
	var dbUser m.User
	user := m.User{
		Username: "testuser",
		Email:    "testemail@email.com",
		Password: "testpassword",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
	req.Header = suite.baseHeader

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	suite.Equal(http.StatusCreated, w.Code)
	suite.Nil(err)
	suite.Equal(user.Username, createdUser.Username)
	suite.Equal(user.Email, createdUser.Email)

	err = suite.db.Where("id = ?", createdUser.ID).First(&dbUser).Error
	suite.Nil(err)
	suite.Equal(user.Username, dbUser.Username)
	suite.Equal(user.Email, dbUser.Email)
	suite.NotNil(dbUser.Password)
	fmt.Println(dbUser.Password)
	suite.NotEqual(user.Password, dbUser.Password)
}

func (suite *UserCreateTestSuite) TestFailedUserCreationBadRequest() {
	//TOOD: Add more test cases and error responses when it will be better handled
	testCases := []m.User{
		{
			Username: "testuser",
			Email:    "testemail@email.com",
			Password: "short",
		},
		{
			Username: "testuser",
			Email:    "invalidemail",
			Password: "testpassword",
		},
	}

	for _, user := range testCases {
		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
		req.Header = suite.baseHeader

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusBadRequest, w.Code)
	}
}

func (suite *UserCreateTestSuite) TestFailedUserCreationConflict() {
	user := m.User{
		Username: "testuser",
		Email:    "testemail@email.com",
		Password: "testpassword",
	}

	suite.db.Model(&m.User{}).Create(&user)

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
	req.Header = suite.baseHeader

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusConflict, w.Code)

}

// func (suite *UserControllerTestSuite) TestGetUser(t *testing.T) {
// 	// Get an existing user
// 	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
// 	req.Header.Set("Authorization", "Bearer "+suite.token)

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var user m.User
// 	err := json.Unmarshal(w.Body.Bytes(), &user)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "authuser", user.Username)
// 	assert.Equal(t, "authuser@email.com", user.Email)
// 	fmt.Println("TestGetUser")
// }

// func (suite *UserControllerTestSuite) TestUpdateUser(t *testing.T) {
// 	// Update an existing user
// 	updatedData := map[string]string{
// 		"username": "updateduser",
// 		"email":    "updatedemail@email.com",
// 	}
// 	jsonData, _ := json.Marshal(updatedData)
// 	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+suite.token)

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var updatedUser m.User
// 	err := json.Unmarshal(w.Body.Bytes(), &updatedUser)
// 	assert.Nil(t, err)
// 	assert.Equal(t, updatedData["username"], updatedUser.Username)
// 	assert.Equal(t, updatedData["email"], updatedUser.Email)
// 	fmt.Println("TestUpdateUser")
// }

// func (suite *UserControllerTestSuite) TestDeleteUser(t *testing.T) {
// 	// Delete an existing user
// 	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
// 	req.Header.Set("Authorization", "Bearer "+suite.token)

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNoContent, w.Code)
// 	fmt.Println("TestDeleteUser")
// }

func (suite *UserCreateTestSuite) TearDownTest() {
	db.Exec("DELETE FROM users WHERE username != 'admin'")
}

func (suite *UserCreateTestSuite) TearDownSuite() {
	suite.pgContainter.Terminate(suite.ctx)
}

func TestUserCreateTestSuite(t *testing.T) {
	suite.Run(t, new(UserCreateTestSuite))
}
