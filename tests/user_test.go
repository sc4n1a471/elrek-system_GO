package tests

import (
	"bytes"
	"elrek-system_GO/api"
	"elrek-system_GO/controllers"
	"elrek-system_GO/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var adminCookies []*http.Cookie
var nonAdminCookies []*http.Cookie
var router *gin.Engine
var userId openapitypes.UUID

func TestSetupUserTest(t *testing.T) {
	dbError := controllers.SetupDB()
	if dbError != nil {
		t.Errorf("Error setting up database: %s", dbError)
	}

	router = api.SetupRouter()
}

func TestLoginWrongPassword(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user1@example.com",
		Password: "stringg",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Wrong password",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 401, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestLoginWrongEmail(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user11@example.com",
		Password: "stringg",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Could not get user: record not found",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 500, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestLoginAsAdmin(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user1@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.UserLoginResponse{}
	correctResponseBody := models.UserLoginResponse{
		Email:   "user1@example.com",
		Id:      openapitypes.UUID(uuid.MustParse("85ee6a8a-3fb8-4a87-9d76-3656524697fb")),
		Name:    "",
		IsAdmin: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)
	adminCookies = w.Result().Cookies()

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: Create user ===================
func TestUserCreate(t *testing.T) {
	requestBody := models.UserCreate{
		Email:    "user5@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 201, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserCreateDuplicate(t *testing.T) {
	requestBody := models.UserCreate{
		Email:    "user1@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User with this email already exists",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 400, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserCreateWithoutCookies(t *testing.T) {
	requestBody := models.UserCreate{
		Email:    "user1@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Not logged in",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 401, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestLoginAsNonAdmin(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user5@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.UserLoginResponse{}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)
	nonAdminCookies = w.Result().Cookies()

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	userId = responseBody.Id
}

func TestUserCreateWithoutAdmin(t *testing.T) {
	requestBody := models.UserCreate{
		Email:    "user5@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 403, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserGetCreatedUser(t *testing.T) {
	var responseBody []models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert2.Contains(t, responseBody, correctResponseBody)
}

// MARK: Update user ===================
// MARK: Update user without admin =====
func TestUserUpdateNameWithoutAdmin(t *testing.T) {
	requestBody := models.UserUpdate{}
	username := "User 5"
	requestBody.Name = &username
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+userId.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserUpdateNameWithoutAdminCheck(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 5",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId.String(), nil)
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

func TestUserUpdateIsAdminWithoutAdmin(t *testing.T) {
	requestBody := models.UserUpdate{}
	isAdmin := true
	requestBody.IsAdmin = &isAdmin
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+userId.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 403, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserUpdateIsAdminWithoutAdminCheck(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 5",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId.String(), nil)
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

// MARK: Update user with admin =====
func TestUserUpdateName(t *testing.T) {
	requestBody := models.UserUpdate{}
	username := "User 55"
	requestBody.Name = &username
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+userId.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserUpdateNameCheck(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 55",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

func TestUserUpdateIsAdmin(t *testing.T) {
	requestBody := models.UserUpdate{}
	isAdmin := true
	requestBody.IsAdmin = &isAdmin
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+userId.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserUpdateIsAdminCheck(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 55",
		IsAdmin:  true,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

// MARK: Delete user ===================
func TestUserDelete(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/"+userId.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUserDeleteGetUsers(t *testing.T) {
	var responseBody []models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 5",
		IsAdmin:  true,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert2.NotContains(t, responseBody, correctResponseBody)
}

func TestUserDeleteGetUser(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		Id:       userId,
		Name:     "User 55",
		IsAdmin:  true,
		IsActive: false,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

func TestUserDeletePermanently(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was permanently deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/permanently/"+userId.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}