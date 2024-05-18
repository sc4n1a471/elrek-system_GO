package tests

import (
	"elrek-system_GO/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	openapitypes "github.com/oapi-codegen/runtime/types"
	assert2 "github.com/stretchr/testify/assert"
)

// MARK: Delete services ===================
func TestServiceDeleteWithoutAdmin(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Access denied"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/services/"+serviceObject.ID.String(), nil)
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestServiceDelete(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was deleted successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/services/"+serviceObject.ID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceDeleteCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive:      false,
		Name:          serviceName + "_Updated",
		UserID:        adminUserID,
		PrevServiceID: openapitypes.UUID{},
		Price:         5001,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == serviceName {
			if service.Price == correctResponseBody.Price &&
				!service.IsActive &&
				service.UserID == adminUserID &&
				service.PrevServiceID == correctResponseBody.PrevServiceID {

				assert.NotEqual(t, correctResponseBody.Name, service.Name)
				return
			}
		}
	}
}

func TestServiceWDPDelete(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was deleted successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/services/"+serviceWDPObject.ID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceWDPDeleteCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive:      false,
		Name:          serviceName + "_DP" + "_Updated",
		UserID:        adminUserID,
		PrevServiceID: openapitypes.UUID{},
		Price:         5001,
		DynamicPrices: &[]models.DynamicPrice{
			{
				Attendees: 3,
				Price:     6001,
			},
			{
				Attendees: 2,
				Price:     7001,
			},
			{
				Attendees: 1,
				Price:     8001,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == serviceName {
			if service.Price == correctResponseBody.Price &&
				!service.IsActive &&
				service.UserID == adminUserID &&
				service.DynamicPrices == correctResponseBody.DynamicPrices {

				assert.NotEqual(t, correctResponseBody.Name, service.Name)
				return
			}
		}
	}
}

//func TestServiceMarkingAsDone(t *testing.T) {
//	doneServiceName := serviceName + "_DONE"
//	doneServicePrice := 5000
//	requestBody := models.ServiceUpdate{
//		Name:  &doneServiceName,
//		Price: &doneServicePrice,
//	}
//	marshalledRequestBody, _ := json.Marshal(requestBody)
//
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("PATCH", "/services/"+serviceObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}

// MARK: Delete user ===================
func TestUserDelete(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/"+halfAdminUserID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

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
		ID:       halfAdminUserID,
		Name:     "User 5",
		IsAdmin:  true,
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
	assert2.NotContains(t, responseBody, correctResponseBody)
}

func TestUserDeleteGetUser(t *testing.T) {
	var responseBody models.UserResponse
	correctResponseBody := models.UserResponse{
		Email:    "user5@example.com",
		ID:       halfAdminUserID,
		Name:     "User 55",
		IsAdmin:  true,
		IsActive: false,
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

func TestUserDeletePermanently(t *testing.T) {
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "User was permanently deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/permanently/"+halfAdminUserID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
