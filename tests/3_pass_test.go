package tests

import (
	"bytes"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var passName string
var service2ID string
var service2Object models.Service

func TestPassSetup(t *testing.T) {
	passName = fmt.Sprint("Pass_", randomID)
	fmt.Println("Pass name: ", passName)
}

func TestPassCreateWithoutLoggingIn(t *testing.T) {
	// Create a pass without logging in
	requestBody := models.PassCreate{
		Name:          "",
		OccasionLimit: nil,
		UserID:        nonAdminUserID,
		Price:         0,
		ServiceIDs:    nil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Not logged in",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes", bytes.NewBuffer(marshalledRequestBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassCreateWithoutAdmin(t *testing.T) {
	// Create a pass without being an admin
	requestBody := models.PassCreate{
		Name:          "",
		OccasionLimit: nil,
		UserID:        nonAdminUserID,
		Price:         0,
		ServiceIDs:    nil,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestPassCreate(t *testing.T) {
	requestBody := models.PassCreate{
		Name:          passName,
		OccasionLimit: nil,
		UserID:        adminUserID,
		Price:         5000,
		ServiceIDs:    []openapitypes.UUID{serviceID},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass was created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/passes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassCreateCheck(t *testing.T) {
	var responseBody []models.Pass
	correctResponseBody := models.Pass{
		IsActive:      true,
		Comment:       nil,
		Duration:      nil,
		Name:          passName,
		OccasionLimit: nil,
		UserID:        adminUserID,
		PrevPassID:    openapitypes.UUID{},
		Price:         5000,
		Services: []models.Service{
			serviceObject,
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	passObject = checkPassEqual(t, responseBody, correctResponseBody, true)
}

func checkPassEqual(t *testing.T, responseBody []models.Pass, correctResponseBody models.Pass, checkServices bool) models.Pass {
	for _, pass := range responseBody {
		if pass.Name == correctResponseBody.Name {
			assert.Equal(t, correctResponseBody.IsActive, pass.IsActive)
			assert.Equal(t, correctResponseBody.Comment, pass.Comment)
			assert.Equal(t, correctResponseBody.Duration, pass.Duration)
			assert.Equal(t, correctResponseBody.OccasionLimit, pass.OccasionLimit)
			assert.Equal(t, correctResponseBody.UserID, pass.UserID)
			//assert.Equal(t, correctResponseBody.PrevPassID, pass.PrevPassID)
			assert.Equal(t, correctResponseBody.Price, pass.Price)
			assert.Equal(t, correctResponseBody.Name, pass.Name)

			if checkServices {
				// Required because getPasses does not return dynamic prices with services
				correctResponseBody.Services[0].DynamicPrices = nil
				checkServicesEqual(t, pass.Services, correctResponseBody.Services)
			}
			return pass
		}
	}
	t.Errorf("Error: Pass not found")
	return models.Pass{}
}
func checkServicesEqual(t *testing.T, responseBody []models.Service, correctResponseBody []models.Service) {
	for _, service := range responseBody {
		for _, correctService := range correctResponseBody {
			if service.Name == correctService.Name {
				assert.Equal(t, correctService.IsActive, service.IsActive)
				assert.Equal(t, correctService.Comment, service.Comment)
				assert.Equal(t, correctService.UserID, service.UserID)
				//assert.Equal(t, correctService.PrevServiceID, service.PrevServiceID)
				assert.Equal(t, correctService.Price, service.Price)
				assert.Equal(t, correctService.Name, service.Name)
				break
			}
		}
	}
}

func TestPassUpdate1(t *testing.T) {
	updatedName := passName + "_Updated"
	updatedComment := "Updated comment"
	requestBody := models.PassUpdate{
		Name:    &updatedName,
		Comment: &updatedComment,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprint("/passes/", passObject.ID), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassUpdate1Check(t *testing.T) {
	var responseBody []models.Pass
	correctComment := "Updated comment"
	correctResponseBody := models.Pass{
		IsActive:      true,
		Comment:       &correctComment,
		Duration:      nil,
		Name:          passName + "_Updated",
		OccasionLimit: nil,
		UserID:        adminUserID,
		PrevPassID:    openapitypes.UUID{},
		Price:         5000,
		Services: []models.Service{
			serviceObject,
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	checkPassEqual(t, responseBody, correctResponseBody, true)
}

func TestPassUpdate2(t *testing.T) {
	updatedDuration := "2_week"
	updatedOccasionLimit := 8
	requestBody := models.PassUpdate{
		Duration:      &updatedDuration,
		OccasionLimit: &updatedOccasionLimit,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprint("/passes/", passObject.ID), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassUpdate2Check(t *testing.T) {
	var responseBody []models.Pass
	correctComment := "Updated comment"
	correctDuration := "2_week"
	correctOccasionLimit := 8
	correctResponseBody := models.Pass{
		IsActive:      true,
		Comment:       &correctComment,
		Duration:      &correctDuration,
		Name:          passName + "_Updated",
		OccasionLimit: &correctOccasionLimit,
		UserID:        adminUserID,
		Price:         5000,
		Services: []models.Service{
			serviceObject,
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	checkPassEqual(t, responseBody, correctResponseBody, true)
}

// createService3 is used to create a second service to add to pass' services
func createService3() string {
	requestBody := models.ServiceCreate{
		Name:  serviceName + "_2",
		Price: 6000,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive:      true,
		Name:          serviceName + "_2",
		UserID:        adminUserID,
		PrevServiceID: openapitypes.UUID{},
		Price:         6000,
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		return "Error: Could not unmarshal response body"
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == correctResponseBody.UserID &&
				service.PrevServiceID == correctResponseBody.PrevServiceID {
				service.DynamicPrices = nil
				service2Object = service
				return service.ID.String()
			}
			return "Error: \"Service attributes do not match\""
		}
	}
	return "Error: Service not found"
}
func TestPassUpdate3(t *testing.T) {
	service2ID = createService3()
	if strings.Contains(service2ID, "Error") {
		t.Errorf(service2ID)
		return
	}

	updatedServiceIDs := []openapitypes.UUID{serviceID, uuid.MustParse(service2ID)}
	requestBody := models.PassUpdate{
		ServiceIDs: &updatedServiceIDs,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprint("/passes/", passObject.ID), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassUpdate3Check(t *testing.T) {
	var responseBody []models.Pass
	correctComment := "Updated comment"
	correctDuration := "2_week"
	correctOccasionLimit := 8
	correctResponseBody := models.Pass{
		IsActive:      true,
		Comment:       &correctComment,
		Duration:      &correctDuration,
		Name:          passName + "_Updated",
		OccasionLimit: &correctOccasionLimit,
		UserID:        adminUserID,
		PrevPassID:    openapitypes.UUID{},
		Price:         5000,
		Services: []models.Service{
			service2Object,
			serviceObject,
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	checkPassEqual(t, responseBody, correctResponseBody, true)
}

func TestPassDelete(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Pass was deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprint("/passes/", passObject.ID), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestPassDeleteCheck(t *testing.T) {
	var responseBody []models.Pass

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/passes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	for _, pass := range responseBody {
		if pass.ID == passObject.ID {
			fmt.Println(pass)
			if pass.IsActive {
				t.Errorf("Error: Pass was not deleted")
			}
			return
		}
	}
}
