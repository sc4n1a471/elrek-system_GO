package tests

import (
	"bytes"
	"elrek-system_GO/api"
	"elrek-system_GO/controllers"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	assert2 "github.com/stretchr/testify/assert"
)

var adminCookies []*http.Cookie
var halfAdminCookies []*http.Cookie
var nonAdminCookies []*http.Cookie
var router *gin.Engine
var halfAdminUserID openapitypes.UUID
var nonAdminUserID openapitypes.UUID
var adminUserID = openapitypes.UUID(uuid.MustParse("16cf214e-a31c-4fb9-b381-2e30df8cc946"))

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

	assert.Equal(t, 401, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestLoginWrongEmail(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user-1@example.com",
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

	assert.Equal(t, http.StatusInternalServerError, w.Code)

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
		ID:      adminUserID,
		Name:    "",
		IsAdmin: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)
	adminCookies = w.Result().Cookies()

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestLoginAsNonAdmin(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user2@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.UserLoginResponse{}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)
	nonAdminCookies = w.Result().Cookies()

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	nonAdminUserID = responseBody.ID
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

	fmt.Println(w.Body.String())
	assert.Equal(t, 201, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: TestUserCreateDuplicate
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

	assert.Equal(t, 500, w.Code) // TODO: Change to 400 when issue #3 is done

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: TestUserCreateWithoutCookies
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

	assert.Equal(t, 401, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: TestUserCreateWithoutAdmin
func TestLoginAsHalfAdmin(t *testing.T) {
	requestBody := models.UserLogin{
		Email:    "user5@example.com",
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.UserLoginResponse{}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(marshalledRequestBody))
	router.ServeHTTP(w, req)
	halfAdminCookies = w.Result().Cookies()

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	halfAdminUserID = responseBody.ID
}

// MARK: TestUserCreateWithoutAdmin
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
	req.AddCookie(halfAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: Get created user
func TestUserGetCreatedUser(t *testing.T) {
	var responseBody []models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		ID:       halfAdminUserID,
		Name:     "",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert2.Contains(t, responseBody, correctResponseBody)
}

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
	req, _ := http.NewRequest("PATCH", "/users/"+halfAdminUserID.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(halfAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: TestUserUpdateNameWithoutAdminCheck
func TestUserUpdateNameWithoutAdminCheck(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		ID:       halfAdminUserID,
		Name:     "User 5",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+halfAdminUserID.String(), nil)
	req.AddCookie(halfAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}

// MARK: UserUpdateIsAdminWithoutAdmin
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
	req, _ := http.NewRequest("PATCH", "/users/"+halfAdminUserID.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(halfAdminCookies[0])
	router.ServeHTTP(w, req)

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
		ID:       halfAdminUserID,
		Name:     "User 5",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+halfAdminUserID.String(), nil)
	req.AddCookie(halfAdminCookies[0])
	router.ServeHTTP(w, req)

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
	req, _ := http.NewRequest("PATCH", "/users/"+halfAdminUserID.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

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
		ID:       halfAdminUserID,
		Name:     "User 55",
		IsAdmin:  false,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+halfAdminUserID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

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

	fmt.Println(nonAdminUserID)
	fmt.Println(halfAdminUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+halfAdminUserID.String(), bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

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
		ID:       halfAdminUserID,
		Name:     "User 55",
		IsAdmin:  true,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+halfAdminUserID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}
