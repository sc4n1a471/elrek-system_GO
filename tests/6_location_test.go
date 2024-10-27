package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"elrek-system_GO/models"
)

var locationName string
var locationObject models.Location

func TestLocationSetup(t *testing.T) {
	locationName = fmt.Sprint("Locarion_", randomID)
	fmt.Println("Location name:", locationName)
}

func checkLocationEqual(t *testing.T, correct models.Location, response []models.Location, IDshouldEqual bool) models.Location {
	for _, location := range response {
		if location.Name == correct.Name {
			assert.Equal(t, correct.Name, location.Name)
			assert.Equal(t, correct.IsActive, location.IsActive)
			if IDshouldEqual {
				assert.Equal(t, correct.ID, location.ID)
			}
			return location
		}
	}
	t.Errorf("Location not found: %v", correct)
	return models.Location{}
}

// MARK: Create location as non-admin
func TestCreateLocationWithoutAdmin(t *testing.T) {
	// Create a location without being an admin
	requestBody := models.LocationCreate{
		Name:     locationName,
		IsActive: true,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/locations", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: Create location as admin
func TestCreateLocationWithAdmin(t *testing.T) {
	// Create a location as an admin
	requestBody := models.LocationCreate{
		Name:     locationName,
		IsActive: true,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Location created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/locations", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestCreateLocationWithAdminCheck(t *testing.T) {
	// Check if the location was created
	var responseBody []models.Location
	correctResponseBody := models.Location{
		Name:     locationName,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	locationObject = checkLocationEqual(t, correctResponseBody, responseBody, false)
}

// MARK: Update location as non-admin
func TestUpdateLocationWithoutAdmin(t *testing.T) {
	// Update the location without being an admin
	requestBody := models.LocationUpdate{
		Name:       &locationName,
		IsActive:   func(b bool) *bool { return &b }(false),
		UpdateOnly: true,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/locations/"+locationObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: Update location as admin 1
func TestUpdateLocationWithAdmin1(t *testing.T) {
	// Update the location as an admin
	updatedName := locationName + "_updated_1"
	requestBody := models.LocationUpdate{
		Name:       &updatedName,
		UpdateOnly: true,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Location was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/locations/"+locationObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
	if correctResponseBody != responseBody {
		t.Errorf("Response body: %v", responseBody)
	}
}

func TestUpdateLocationCheck1(t *testing.T) {
	// Check if the location was updated
	var responseBody []models.Location
	correctResponseBody := models.Location{
		ID:       locationObject.ID,
		Name:     locationName + "_updated_1",
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	locationObject = checkLocationEqual(t, correctResponseBody, responseBody, true)
}

// MARK: Update location as admin 2
func TestUpdateLocationWithAdmin2(t *testing.T) {
	// Update the location as an admin
	updatedName := locationName + "_updated_2"
	requestBody := models.LocationUpdate{
		Name:       &updatedName,
		UpdateOnly: false,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Location was updated successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/locations/"+locationObject.ID.String(), bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestUpdateLocationCheck2(t *testing.T) {
	// Check if the location was updated
	var responseBody []models.Location
	correctResponseBody := models.Location{
		ID:       locationObject.ID,
		Name:     locationName + "_updated_2",
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	locationObject = checkLocationEqual(t, correctResponseBody, responseBody, false)
}

// MARK: Delete location as non-admin
func TestDeleteLocationWithoutAdmin(t *testing.T) {
	// Delete the location without being an admin
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Access denied",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/locations/"+locationObject.ID.String(), nil)
	req.AddCookie(nonAdminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

// MARK: Delete location as admin
func TestDeleteLocationWithAdmin(t *testing.T) {
	// Delete the location as an admin
	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Location deleted successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/locations/"+locationObject.ID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestDeleteLocationCheck1(t *testing.T) {
	// Check if the location was deleted
	var responseBody []models.Location

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	for _, location := range responseBody {
		if location.Name == locationObject.Name {
			if location.IsActive {
				t.Errorf("Location was not deleted: %v", location)
			}
			return
		}
	}
}

func TestDeleteLocationCheck2(t *testing.T) {
	// Check if the location was deleted
	var responseBody models.Location

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations/"+locationObject.ID.String(), nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}

	if responseBody.IsActive {
		t.Errorf("Location was not deleted: %v", responseBody)
	}
}

// MARK: Create location as admin for further testing
func TestCreateLocation(t *testing.T) {
	// Create a location as an admin
	requestBody := models.LocationCreate{
		Name:     locationName,
		IsActive: true,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Location created successfully",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/locations", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}

func TestCreateLocationCheck(t *testing.T) {
	// Check if the location was created
	var responseBody []models.Location
	correctResponseBody := models.Location{
		Name:     locationName,
		IsActive: true,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		t.Errorf("Error while unmarshalling response body: %v", err)
	}
	locationObject = checkLocationEqual(t, correctResponseBody, responseBody, false)
}
