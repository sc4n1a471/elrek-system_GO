package tests

import (
	"bytes"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	openapitypes "github.com/oapi-codegen/runtime/types"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var validFrom = time.Now().Round(time.Second)
var validUntil = time.Now().Add(time.Hour * 24).Round(time.Second)
var occasionLimit = 2
var passInUseObject models.PassInUse
var passInUseWDPObject models.PassInUse
var newPassName string
var passObject2 models.Pass

var updatedValidFrom = validFrom.Add(-(time.Hour * 24)).Round(time.Second)
var updatedValidUntil = validUntil.Add(time.Hour * 24).Round(time.Second)
var updatedComment = "Updated comment"
var updatedOccasions = 1

func checkPassInUseEqual(t *testing.T, responseBody []models.PassInUse, correctResponseBody models.PassInUse) models.PassInUse {
	for _, passInUse := range responseBody {
		if passInUse.PassID == correctResponseBody.PassID &&
			passInUse.IsActive == correctResponseBody.IsActive {

			if passInUse.Comment != nil {
				assert.Equal(t, correctResponseBody.Comment, *passInUse.Comment)
			} else {
				assert.Equal(t, correctResponseBody.Comment, passInUse.Comment)
			}

			assert.Equal(t, correctResponseBody.UserID, passInUse.UserID)
			assert.Equal(t, correctResponseBody.PayerID, passInUse.PayerID)
			assert.Equal(t, correctResponseBody.ValidFrom, passInUse.ValidFrom)
			assert.Equal(t, correctResponseBody.ValidUntil, passInUse.ValidUntil)
			assert.Equal(t, correctResponseBody.Occasions, passInUse.Occasions)

			checkPassEqual(t, []models.Pass{correctResponseBody.Pass}, passInUse.Pass, false, false)

			return passInUse
		}
	}
	t.Error("Pass in use not found")
	return models.PassInUse{}
}

func TestPassInUseSetup(t *testing.T) {
	newPassName = fmt.Sprint(passName+"InUse", randomID)
	fmt.Println("newPassName", newPassName)
}

func TestPassInUseCreateWithoutLoggingIn(t *testing.T) {
	requestBody := models.PassInUseCreate{
		UserID:     openapitypes.UUID{},
		PassID:     openapitypes.UUID{},
		PayerID:    openapitypes.UUID{},
		ValidFrom:  time.Time{},
		ValidUntil: time.Time{},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Not logged in",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassInUseCreateWithoutWithoutAdmin(t *testing.T) {
	requestBody := models.PassInUseCreate{
		UserID:     openapitypes.UUID{},
		PassID:     openapitypes.UUID{},
		PayerID:    openapitypes.UUID{},
		ValidFrom:  time.Time{},
		ValidUntil: time.Time{},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
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
func TestPassInUseCreate(t *testing.T) {
	passObject2 = createPass(t)
	requestBody := models.PassInUseCreate{
		UserID:     adminUserID,
		PassID:     passObject2.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass in use was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassInUseCreateCheck(t *testing.T) {
	var responseBody []models.PassInUse
	correctResponseBody := models.PassInUse{
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
	req, _ := http.NewRequest("GET", "/passes_in_use", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	passInUseObject = checkPassInUseEqual(t, responseBody, correctResponseBody)
}

func TestPassInUseCreateWithInvalidPassID(t *testing.T) {
	requestBody := models.PassInUseCreate{
		UserID:     adminUserID,
		PassID:     openapitypes.UUID{},
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Could not get pass: record not found",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassInUseCreateWithInvalidPayerID(t *testing.T) {
	requestBody := models.PassInUseCreate{
		UserID:     adminUserID,
		PassID:     passObject.ID,
		PayerID:    openapitypes.UUID{},
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Could not get user (payer):",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert2.Contains(t, responseBody.Message, correctResponseBody.Message)
}

func TestPassInUseUpdate(t *testing.T) {
	requestBody := models.PassInUseUpdate{
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
	req, _ := http.NewRequest("PATCH", "/passes_in_use/"+passInUseObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassInUseUpdateCheck(t *testing.T) {
	var responseBody []models.PassInUse
	correctResponseBody := models.PassInUse{
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
	req, _ := http.NewRequest("GET", "/passes_in_use", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	passInUseObject = checkPassInUseEqual(t, responseBody, correctResponseBody)
}

// occasions limit is not reached, is at 1/2
func TestPassInUseValidityCheck1(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes_in_use/"+passInUseObject.ID.String()+"/validity", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Body.String())
}

func TestPassInUseCreateWDP(t *testing.T) {
	requestBody := models.PassInUseCreate{
		UserID:     adminUserID,
		PassID:     passWDPObject.ID,
		PayerID:    nonAdminUserID,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass in use was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassInUseCreateWDPCheck(t *testing.T) {
	var responseBody []models.PassInUse
	correctResponseBody := models.PassInUse{
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
	req, _ := http.NewRequest("GET", "/passes_in_use", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	passInUseWDPObject = checkPassInUseEqual(t, responseBody, correctResponseBody)
}

//// occasions limit is reached, not yet invalidated
//func TestPassInUseUsage(t *testing.T) {
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{
//		Message: "Pass in use was used successfully",
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/passes_in_use/"+passInUseObject.ID.String()+"/use", nil)
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
//func TestPassInUseUsageCheck(t *testing.T) {
//	var responseBody []models.PassInUse
//	correctResponseBody := models.PassInUse{
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
//	req, _ := http.NewRequest("GET", "/passes_in_use", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//
//	passInUseObject = checkPassInUseEqual(t, responseBody, correctResponseBody)
//}
//
//// Try to use a passInUse with limit reached, should invalidate it
//func TestPassInUseUsage2(t *testing.T) {
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{
//		Message: "Pass in use is not valid",
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/passes_in_use/"+passInUseObject.ID.String()+"/use", nil)
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
//func TestPassInUseValidityCheck2(t *testing.T) {
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/passes_in_use/"+passInUseObject.ID.String()+"/validity", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, "false", w.Body.String())
//}
//
//// Test valid_until validation
//func TestPassInUseCreate2(t *testing.T) {
//	passObject = createPass(t)
//	requestBody := models.PassInUseCreate{
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
//	req, _ := http.NewRequest("POST", "/passes_in_use", bytes.NewBuffer(marshalledRequestBody))
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
//func TestPassInUseValidityCheck3(t *testing.T) {
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/passes_in_use/"+passInUseObject.ID.String()+"/validity", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, "false", w.Body.String())
//}
