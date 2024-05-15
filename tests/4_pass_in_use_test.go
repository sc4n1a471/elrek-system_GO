package tests

import (
	"bytes"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	openapitypes "github.com/oapi-codegen/runtime/types"
	assert2 "github.com/stretchr/testify/assert"
)

var validFrom = time.Now().Round(time.Second)
var validUntil = time.Now().Add(time.Hour * 24).Round(time.Second)
var occasionLimit = 2
var activePassObject models.ActivePass
var activePassWDPObject models.ActivePass
var newPassName string
var passObject2 models.Pass

var updatedValidFrom = validFrom.Add(-(time.Hour * 24)).Round(time.Second)
var updatedValidUntil = validUntil.Add(time.Hour * 24).Round(time.Second)
var updatedComment = "Updated comment"
var updatedOccasions = 1

func checkactivePassEqual(t *testing.T, responseBody []models.ActivePass, correctResponseBody models.ActivePass) models.ActivePass {
	for _, activePass := range responseBody {
		if activePass.PassID == correctResponseBody.PassID &&
			activePass.IsActive == correctResponseBody.IsActive {

			if activePass.Comment != nil {
				assert.Equal(t, correctResponseBody.Comment, *activePass.Comment)
			} else {
				assert.Equal(t, correctResponseBody.Comment, activePass.Comment)
			}

			assert.Equal(t, correctResponseBody.UserID, activePass.UserID)
			assert.Equal(t, correctResponseBody.PayerID, activePass.PayerID)
			assert.Equal(t, correctResponseBody.ValidFrom, activePass.ValidFrom)
			assert.Equal(t, correctResponseBody.ValidUntil, activePass.ValidUntil)
			assert.Equal(t, correctResponseBody.Occasions, activePass.Occasions)

			checkPassEqual(t, []models.Pass{correctResponseBody.Pass}, activePass.Pass, false, false)

			return activePass
		}
	}
	t.Error("Pass in use not found")
	return models.ActivePass{}
}

func TestactivePassSetup(t *testing.T) {
	newPassName = fmt.Sprint(passName+"InUse", randomID)
	fmt.Println("newPassName", newPassName)
}

func TestactivePassCreateWithoutLoggingIn(t *testing.T) {
	requestBody := models.ActivePassCreate{
		UserID:    openapitypes.UUID{},
		PassID:    openapitypes.UUID{},
		PayerID:   openapitypes.UUID{},
		ValidFrom: time.Time{},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Not logged in",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestactivePassCreateWithoutWithoutAdmin(t *testing.T) {
	requestBody := models.ActivePassCreate{
		UserID:    openapitypes.UUID{},
		PassID:    openapitypes.UUID{},
		PayerID:   openapitypes.UUID{},
		ValidFrom: time.Time{},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// createPass is used to create a pass with basic serviceObject
func createPass(t *testing.T) models.Pass {
	requestBody := models.PassCreate{
		Name:          newPassName,
		OccasionLimit: &occasionLimit,
		UserID:        adminUserID,
		Price:         5000,
		ServiceIDs:    []openapitypes.UUID{serviceObject.ID},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	var responseBody []models.Pass
	correctResponseBody := models.Pass{
		IsActive:      true,
		Comment:       nil,
		Duration:      nil,
		Name:          newPassName,
		OccasionLimit: &occasionLimit,
		UserID:        adminUserID,
		PrevPassID:    openapitypes.UUID{},
		Price:         5000,
		Services: []models.Service{
			serviceObject,
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Error(err)
		return models.Pass{}
	}

	return checkPassEqual(t, responseBody, correctResponseBody, false, false)
}
func TestactivePassCreate(t *testing.T) {
	passObject2 = createPass(t)
	requestBody := models.ActivePassCreate{
		UserID:     adminUserID,
		PassID:     passObject2.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: &validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass in use was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestactivePassCreateCheck(t *testing.T) {
	var responseBody []models.ActivePass
	correctResponseBody := models.ActivePass{
		IsActive:   true,
		Comment:    nil,
		UserID:     adminUserID,
		PassID:     passObject2.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  &validFrom,
		ValidUntil: &validUntil,
		Occasions:  0,
		Pass:       passObject2,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/active-passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	activePassObject = checkactivePassEqual(t, responseBody, correctResponseBody)
}

func TestactivePassCreateWithInvalidPassID(t *testing.T) {
	requestBody := models.ActivePassCreate{
		UserID:     adminUserID,
		PassID:     openapitypes.UUID{},
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: &validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Could not get pass: record not found",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestactivePassCreateWithInvalidPayerID(t *testing.T) {
	requestBody := models.ActivePassCreate{
		UserID:     adminUserID,
		PassID:     passObject.ID,
		PayerID:    openapitypes.UUID{},
		ValidFrom:  validFrom,
		ValidUntil: &validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Could not get user (payer):",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert2.Contains(t, responseBody.Message, correctResponseBody.Message)
}

func TestactivePassUpdate(t *testing.T) {
	requestBody := models.ActivePassUpdate{
		Comment:    &updatedComment,
		Occasions:  &updatedOccasions,
		ValidFrom:  &updatedValidFrom,
		ValidUntil: &updatedValidUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass in use was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/active-passes/"+activePassObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestactivePassUpdateCheck(t *testing.T) {
	var responseBody []models.ActivePass
	correctResponseBody := models.ActivePass{
		IsActive:   true,
		Comment:    &updatedComment,
		UserID:     adminUserID,
		PassID:     passObject2.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  &updatedValidFrom,
		ValidUntil: &updatedValidUntil,
		Occasions:  updatedOccasions,
		Pass:       passObject2,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/active-passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	activePassObject = checkactivePassEqual(t, responseBody, correctResponseBody)
}

// occasions limit is not reached, is at 1/2
func TestactivePassValidityCheck1(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/active-passes/"+activePassObject.ID.String()+"/validity", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Body.String())
}

func TestactivePassCreateWDP(t *testing.T) {
	requestBody := models.ActivePassCreate{
		UserID:     adminUserID,
		PassID:     passWDPObject.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: &validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass in use was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestactivePassCreateWDPCheck(t *testing.T) {
	var responseBody []models.ActivePass
	correctResponseBody := models.ActivePass{
		IsActive:   true,
		Comment:    nil,
		UserID:     adminUserID,
		PassID:     passWDPObject.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  &validFrom,
		ValidUntil: &validUntil,
		Occasions:  0,
		Pass:       passWDPObject,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/active-passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	activePassWDPObject = checkactivePassEqual(t, responseBody, correctResponseBody)
}

//// occasions limit is reached, not yet invalidated
//func TestactivePassUsage(t *testing.T) {
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{
//		Message: "Pass in use was used successfully",
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/active-passes/"+activePassObject.ID.String()+"/use", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}
//func TestactivePassUsageCheck(t *testing.T) {
//	var responseBody []models.ActivePass
//	correctResponseBody := models.ActivePass{
//		IsActive:   true,
//		Comment:    &updatedComment,
//		UserID:     adminUserID,
//		PassID:     passObject.ID,
//		PayerID:    nonAdminUserID,
//		ValidFrom:  &updatedValidFrom,
//		ValidUntil: &updatedValidUntil,
//		Occasions:  updatedOccasions + 1,
//		Pass:       passObject,
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/active-passes", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//
//	activePassObject = checkactivePassEqual(t, responseBody, correctResponseBody)
//}
//
//// Try to use a activePass with limit reached, should invalidate it
//func TestactivePassUsage2(t *testing.T) {
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{
//		Message: "Pass in use is not valid",
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/active-passes/"+activePassObject.ID.String()+"/use", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}
//func TestactivePassValidityCheck2(t *testing.T) {
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/active-passes/"+activePassObject.ID.String()+"/validity", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, "false", w.Body.String())
//}
//
//// Test valid_until validation
//func TestactivePassCreate2(t *testing.T) {
//	passObject = createPass(t)
//	requestBody := models.ActivePassCreate{
//		UserID:     adminUserID,
//		PassID:     passObject.ID,
//		PayerID:    nonAdminUserID,
//		ValidFrom:  validFrom.Add(-(time.Hour * 24)),
//		ValidUntil: validFrom,
//	}
//	marshalledRequestBody, _ := json.Marshal(requestBody)
//
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{
//		Message: "Pass in use was created successfully",
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("POST", "/active-passes", bytes.NewBuffer(marshalledRequestBody))
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusCreated, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}
//func TestactivePassValidityCheck3(t *testing.T) {
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/active-passes/"+activePassObject.ID.String()+"/validity", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, "false", w.Body.String())
//}
