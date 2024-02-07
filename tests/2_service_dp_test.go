package tests

import (
	"bytes"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

var serviceID openapitypes.UUID
var serviceWDPID openapitypes.UUID
var serviceName string

func TestServiceSetup(t *testing.T) {
	randomID := rand.Intn(1000)
	serviceName = fmt.Sprint("Service_", randomID)
	fmt.Println("Service name: ", serviceName)
}

func TestServiceCreateWithoutLoggingIn(t *testing.T) {
	requestBody := models.ServiceCreate{
		Name:  serviceName,
		Price: 5000,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Not logged in"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(marshalledRequestBody))
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceCreateWithoutAdmin(t *testing.T) {
	requestBody := models.ServiceCreate{
		Name:  serviceName,
		Price: 5000,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Access denied"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(marshalledRequestBody))

	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestServiceCreate(t *testing.T) {
	requestBody := models.ServiceCreate{
		Name:  serviceName,
		Price: 5000,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was created successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceCreateCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive:      true,
		Name:          serviceName,
		UserID:        adminUserID,
		PrevServiceID: openapitypes.UUID{},
		Price:         5000,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == correctResponseBody.UserID &&
				service.PrevServiceID == correctResponseBody.PrevServiceID {

				serviceID = service.ID
				assert.Equal(t, correctResponseBody.Name, service.Name)
				return
			}
			fmt.Println(service.Price, correctResponseBody.Price)
			fmt.Println(service.IsActive, correctResponseBody.IsActive)
			fmt.Println(service.UserID, correctResponseBody.UserID)
			fmt.Println(service.PrevServiceID, correctResponseBody.PrevServiceID)
			t.Errorf("Service attributes do not match")
			return
		}
	}
	t.Errorf("Service not found")
}

func TestServiceCreateWithDP(t *testing.T) {
	requestBody := models.ServiceCreate{
		Name:  serviceName + "_DP",
		Price: 5000,
		DynamicPriceCreateUpdate: &[]models.DynamicPriceCreateUpdate{
			{
				Attendees: 3,
				Price:     6000,
			},
			{
				Attendees: 2,
				Price:     7000,
			},
			{
				Attendees: 1,
				Price:     8000,
			},
		},
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was created successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, responseBody, correctResponseBody)
}
func TestServiceCreateWithDPCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive:      true,
		Name:          serviceName + "_DP",
		UserID:        adminUserID,
		PrevServiceID: openapitypes.UUID{},
		Price:         5000,
		DynamicPrices: &[]models.DynamicPrice{
			{
				Attendees: 3,
				Price:     6000,
			},
			{
				Attendees: 2,
				Price:     7000,
			},
			{
				Attendees: 1,
				Price:     8000,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == serviceName+"_DP" {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID &&
				service.PrevServiceID == correctResponseBody.PrevServiceID {

				serviceWDPID = service.ID
				assert.Equal(t, correctResponseBody.Name, service.Name)
			}
		}
	}
}

func TestServiceUpdateName(t *testing.T) {
	updatedName := serviceName + "_Updated"
	requestBody := models.ServiceUpdate{
		Name: &updatedName,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/services/"+serviceID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceUpdateNameCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		Name:   serviceName + "_Updated",
		UserID: adminUserID,
		Price:  5000,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID {

				serviceID = service.ID
				assert.Equal(t, correctResponseBody.Name, service.Name)
				return
			}
			t.Errorf("Service attributes do not match")
		}
	}
	t.Errorf("Service not found")
}

func TestServiceUpdatePrice(t *testing.T) {
	updatedPrice := 5001
	requestBody := models.ServiceUpdate{
		Price: &updatedPrice,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/services/"+serviceID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceUpdatePriceCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		Name:   serviceName + "_Updated",
		UserID: adminUserID,
		Price:  5001,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID {

				serviceID = service.ID
				assert.Equal(t, correctResponseBody.Name, service.Name)
				return
			}
		}
	}
	t.Errorf("Service not found")
}

//func TestServiceUpdateRevert(t *testing.T) {
//	requestBody := models.ServiceCreate{
//		Name:  serviceName,
//		Price: 5000,
//	}
//	marshalledRequestBody, _ := json.Marshal(requestBody)
//
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("PATCH", "/services/"+serviceID.String(), bytes.NewBuffer(marshalledRequestBody))
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	// MARK: Asserts ================
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}
//func TestServiceUpdateRevertCheck(t *testing.T) {
//	responseBody := models.ServiceList{}
//	correctResponseBody := models.Service{
//		Name:    serviceName,
//		PayerID: adminUserID,
//		Price:   5000,
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/services", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	// MARK: Asserts ================
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//	}
//
//	for _, service := range responseBody {
//		if service.Name == correctResponseBody.Name {
//			if service.Price == correctResponseBody.Price &&
//				service.IsActive &&
//				service.PayerID == adminUserID {
//
//				serviceID = service.ID
//				assert.Equal(t, correctResponseBody.Name, service.Name)
//				return
//			}
//		}
//	}
//	t.Errorf("Service not found")
//}

func TestServiceWDPUpdateName(t *testing.T) {
	updatedName := serviceName + "_DP" + "_Updated"
	requestBody := models.ServiceUpdate{
		Name: &updatedName,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/services/"+serviceWDPID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceWDPUpdateNameCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive: true,
		Name:     serviceName + "_DP" + "_Updated",
		UserID:   adminUserID,
		Price:    5000,
		DynamicPrices: &[]models.DynamicPrice{
			{
				Attendees: 3,
				Price:     6000,
			},
			{
				Attendees: 2,
				Price:     7000,
			},
			{
				Attendees: 1,
				Price:     8000,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID {

				if models.AreDPsEqualInAttPri(*service.DynamicPrices, *correctResponseBody.DynamicPrices) {
					serviceWDPID = service.ID
					assert.Equal(t, correctResponseBody.Name, service.Name)
					return
				}

				fmt.Println(service.DynamicPrices, correctResponseBody.DynamicPrices)
				t.Errorf("Dynamic prices does not match")
				return
			}
			t.Errorf("Service attributes do not match")
			return
		}
	}
	t.Errorf("Service not found")
}

func TestServiceWDPUpdatePrice(t *testing.T) {
	updatedPrice := 5001
	requestBody := models.ServiceUpdate{
		Price: &updatedPrice,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/services/"+serviceWDPID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceWDPUpdatePriceCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive: true,
		Name:     serviceName + "_DP" + "_Updated",
		UserID:   adminUserID,
		Price:    5001,
		DynamicPrices: &[]models.DynamicPrice{
			{
				Attendees: 3,
				Price:     6000,
			},
			{
				Attendees: 2,
				Price:     7000,
			},
			{
				Attendees: 1,
				Price:     8000,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/services", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID {

				if models.AreDPsEqualInAttPri(*service.DynamicPrices, *correctResponseBody.DynamicPrices) {
					serviceWDPID = service.ID
					assert.Equal(t, correctResponseBody.Name, service.Name)
					return
				}

				fmt.Println(service.DynamicPrices, correctResponseBody.DynamicPrices)
				t.Errorf("Dynamic prices does not match")
				return
			}
			t.Errorf("Service attributes do not match")
			return
		}
	}
	t.Errorf("Service not found")
}

func TestServiceWDPUpdateDP(t *testing.T) {
	requestBody := models.ServiceUpdate{
		DynamicPrices: &[]models.DynamicPriceCreateUpdate{
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
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/services/"+serviceWDPID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestServiceWDPUpdateDPCheck(t *testing.T) {
	responseBody := models.ServiceList{}
	correctResponseBody := models.Service{
		IsActive: true,
		Name:     serviceName + "_DP" + "_Updated",
		UserID:   adminUserID,
		Price:    5001,
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

	// MARK: Asserts ================
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	for _, service := range responseBody {
		if service.Name == correctResponseBody.Name {
			if service.Price == correctResponseBody.Price &&
				service.IsActive &&
				service.UserID == adminUserID {

				if models.AreDPsEqualInAttPri(*service.DynamicPrices, *correctResponseBody.DynamicPrices) {
					serviceWDPID = service.ID
					assert.Equal(t, correctResponseBody.Name, service.Name)
					return
				}

				fmt.Println(service.DynamicPrices, correctResponseBody.DynamicPrices)
				t.Errorf("Dynamic prices does not match")
				return
			}
			fmt.Println(service, correctResponseBody)
			t.Errorf("Service attributes do not match")
			return
		}
	}
	t.Errorf("Service not found")
}

//func TestServiceWDPUpdateRevert(t *testing.T) {
//	revertedName := serviceName + "_DP"
//	revertedPrice := 5000
//	requestBody := models.ServiceUpdate{
//		Name:  &revertedName,
//		Price: &revertedPrice,
//		DynamicPrices: &[]models.DynamicPriceCreateUpdate{
//			{
//				Attendees: 3,
//				Price:     6000,
//			},
//			{
//				Attendees: 2,
//				Price:     7000,
//			},
//			{
//				Attendees: 1,
//				Price:     8000,
//			},
//		},
//	}
//	marshalledRequestBody, _ := json.Marshal(requestBody)
//
//	responseBody := models.MessageOnlyResponse{}
//	correctResponseBody := models.MessageOnlyResponse{Message: "Service was updated successfully"}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("PATCH", "/services/"+serviceWDPID.String(), bytes.NewBuffer(marshalledRequestBody))
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	// MARK: Asserts ================
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//	}
//	assert.Equal(t, correctResponseBody, responseBody)
//}
//func TestServiceWDPUpdateRevertCheck(t *testing.T) {
//	responseBody := models.ServiceList{}
//	correctResponseBody := models.ServiceWDP{
//		IsActive: true,
//		Name:     serviceName + "_DP",
//		PayerID:  adminUserID,
//		Price:    5000,
//		DynamicPrices: &[]models.DynamicPrice{
//			{
//				Attendees: 3,
//				Price:     6000,
//			},
//			{
//				Attendees: 2,
//				Price:     7000,
//			},
//			{
//				Attendees: 1,
//				Price:     8000,
//			},
//		},
//	}
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/services", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	// MARK: Asserts ================
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//	}
//
//	for _, service := range responseBody {
//		if service.Name == correctResponseBody.Name {
//			if service.Price == correctResponseBody.Price &&
//				service.IsActive &&
//				service.PayerID == adminUserID {
//
//				if models.AreDPsEqualInAttPri(*service.DynamicPrices, *correctResponseBody.DynamicPrices) {
//					serviceWDPID = service.ID
//					assert.Equal(t, correctResponseBody.Name, service.Name)
//					return
//				}
//
//				fmt.Println(service.DynamicPrices, correctResponseBody.DynamicPrices)
//				t.Errorf("Dynamic prices does not match")
//				return
//			}
//			t.Errorf("Service attributes do not match")
//			return
//		}
//	}
//	t.Errorf("Service not found")
//}
