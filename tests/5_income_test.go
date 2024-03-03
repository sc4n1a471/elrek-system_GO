package tests

import (
	"bytes"
	"elrek-system_GO/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var incomeName string
var incomeNameDP string

var incomeNameDPwoPIU string
var incomeNameMU string
var incomeNameMU3 string
var incomeBasic models.Income
var incomeDP models.Income
var createdAt = time.Now().Add(-time.Hour * 24)
var testUsers []models.User
var testUsersIDs []openapitypes.UUID

var testUser1 openapitypes.Email
var testUser2 openapitypes.Email
var testUser3 openapitypes.Email
var testUser4 openapitypes.Email
var testUser5 openapitypes.Email

func checkIncomeEquality(
	t *testing.T,
	expected models.Income,
	actuals []models.Income,
	checkCreatedAt bool,
	wrongAmount bool,
	multipleUsers []models.User) {

	for _, actual := range actuals {
		if actual.Name == nil {
			continue
		}

		if *actual.Name == *expected.Name {
			assert.Equal(t, expected.IsActive, actual.IsActive)

			if wrongAmount {
				assert.NotEqual(t, expected.Amount, actual.Amount)
			} else {
				assert.Equal(t, expected.Amount, actual.Amount)
			}
			assert.Equal(t, expected.Comment, actual.Comment)
			assert.Equal(t, expected.UserID, actual.UserID)

			assert.Equal(t, expected.PassInUseID, actual.PassInUseID)
			if expected.PassInUse != nil {
				assert.Equal(t, expected.PassInUse.ID, actual.PassInUse.ID)
			}

			assert.Equal(t, expected.ServiceID, actual.ServiceID)
			if expected.Service != nil {
				assert.Equal(t, expected.Service.Name, actual.Service.Name)
			}

			if expected.PayerID != actual.PayerID {
				continue
			} else {
			}

			assert.Equal(t, expected.IsPaid, actual.IsPaid)

			if checkCreatedAt {
				assert.Equal(t, expected.CreatedAt.Format(time.RFC3339), actual.CreatedAt.Format(time.RFC3339))
			}

			incomeBasic = actual
			return
		}
	}
	t.Errorf("Income with this name was not found: %s", *expected.Name)
}

func TestIncomeSetup(t *testing.T) {
	incomeName = fmt.Sprint("Income", randomID)
	incomeNameDP = fmt.Sprint("IncomeDP", randomID)
	incomeNameDPwoPIU = fmt.Sprint("IncomeDPwoPIU", randomID)
	incomeNameMU = fmt.Sprint("IncomeMU", randomID)
	incomeNameMU3 = fmt.Sprint("IncomeMU3", randomID)

	testUser1 = openapitypes.Email(fmt.Sprint("user1_", randomID, "@example.com"))
	testUser2 = openapitypes.Email(fmt.Sprint("user2_", randomID, "@example.com"))
	testUser3 = openapitypes.Email(fmt.Sprint("user3_", randomID, "@example.com"))
	testUser4 = openapitypes.Email(fmt.Sprint("user4_", randomID, "@example.com"))
	testUser5 = openapitypes.Email(fmt.Sprint("user5_", randomID, "@example.com"))
}

func TestCreateIncomeBasic(t *testing.T) {
	requestBody := models.IncomeCreate{
		ServiceID: &serviceWoPassObject.ID,
		PayerID:   nonAdminUserID,
		Name:      &incomeName,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Income was created successfully",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/incomes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestCreateIncomeBasicCheck(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      serviceWoPassObject.Price,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWoPassObject.ID,
		Service:     &serviceWoPassObject,
		PayerID:     nonAdminUserID,
		Name:        &incomeName,
		IsPaid:      false,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	checkIncomeEquality(t, correctResponseBody, responseBody, false, false, nil)
}

// Check for income that was created when TestPassInUseCreate was run
//func TestCreateIncomePIUCheck(t *testing.T) {
//	var responseBody []models.Income
//	name := "Bérlet vásárlás"
//	correctResponseBody := models.Income{
//		IsActive:    true,
//		Amount:      passInUseObject.Pass.Price,
//		Comment:     nil,
//		UserID:      adminUserID,
//		PassInUseID: &passInUseObject.ID,
//		PassInUse:   &passInUseObject,
//		ServiceID:   nil,
//		Service:     nil,
//		PayerID:     nonAdminUserID,
//		Name:        &name,
//		IsPaid:      false,
//	}
//
//	w := httptest.NewRecorder()
//	req := httptest.NewRequest("GET", "/incomes", nil)
//	req.AddCookie(adminCookies[0])
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
//	if err != nil {
//		t.Error(err)
//	}
//	checkIncomeEquality(t, correctResponseBody, responseBody, false, false)
//}

// limit is 0 before this, this is an unused PIU
func TestCreateIncomeDP(t *testing.T) {
	requestBody := models.IncomeCreate{
		ServiceID: &serviceWDPObject.ID,
		PayerID:   nonAdminUserID,
		Name:      &incomeNameDP,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Income was created successfully",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/incomes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestCreateIncomeDPCheck(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      0,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWDPObject.ID,
		Service:     &serviceWDPObject,
		PayerID:     nonAdminUserID,
		Name:        &incomeNameDP,
		IsPaid:      false,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	checkIncomeEquality(t, correctResponseBody, responseBody, false, false, nil)
}

func TestCreateIncomeDPwoPIU(t *testing.T) {
	requestBody := models.IncomeCreate{
		ServiceID: &serviceWithDPWoPIUObject.ID,
		PayerID:   nonAdminUserID,
		Name:      &incomeNameDPwoPIU,
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Income was created successfully",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/incomes", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestCreateIncomeDPwoPIUCheck(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      (*serviceWithDPWoPIUObject.DynamicPrices)[2].Price,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWithDPWoPIUObject.ID,
		Service:     &serviceWithDPWoPIUObject,
		PayerID:     nonAdminUserID,
		Name:        &incomeNameDPwoPIU,
		IsPaid:      false,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	checkIncomeEquality(t, correctResponseBody, responseBody, false, false, nil)
}

func createTestUsers() {

	requestBody := models.UserCreate{
		Email:    testUser1,
		Password: "string",
	}
	marshalledRequestBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	requestBody = models.UserCreate{
		Email:    testUser2,
		Password: "string",
	}
	marshalledRequestBody, _ = json.Marshal(requestBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	requestBody = models.UserCreate{
		Email:    testUser3,
		Password: "string",
	}
	marshalledRequestBody, _ = json.Marshal(requestBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	requestBody = models.UserCreate{
		Email:    testUser4,
		Password: "string",
	}
	marshalledRequestBody, _ = json.Marshal(requestBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	requestBody = models.UserCreate{
		Email:    testUser5,
		Password: "string",
	}
	marshalledRequestBody, _ = json.Marshal(requestBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewReader(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)
}
func getTestUsers() {
	var responseBody []models.User

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	err := json.Unmarshal([]byte(w.Body.String()), &responseBody)
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range responseBody {
		if user.Email == testUser1 ||
			user.Email == testUser2 ||
			user.Email == testUser3 ||
			user.Email == testUser4 ||
			user.Email == testUser5 {

			testUsers = append(testUsers, user)
			testUsersIDs = append(testUsersIDs, user.ID)
		}
	}
}

// Checks for 4 users (standard price)
func TestCreateIncomeDPMultipleUsers(t *testing.T) {
	createTestUsers()
	getTestUsers()

	requestBody := models.IncomeCreateMultipleUsers{
		PayerIDs:     testUsersIDs[:4],
		ServiceIDs:   &[]openapitypes.UUID{serviceWDPObject.ID},
		PassInUseIDs: nil,
		Comment:      nil,
		CreatedAt:    &createdAt,
		Name:         &incomeNameMU,
	}

	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Incomes were created successfully",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/incomes/multiple_users", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestCreateIncomeDPMultipleUsersCheck(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      serviceWDPObject.Price,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWDPObject.ID,
		Service:     &serviceWDPObject,
		Name:        &incomeNameMU,
		IsPaid:      false,
		CreatedAt:   createdAt,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	for _, user := range testUsers[:4] {
		correctResponseBody.PayerID = user.ID
		checkIncomeEquality(t, correctResponseBody, responseBody, true, false, testUsers[:4])
	}
}

// Checks for 3 users only
func TestCreateIncomeDPMultipleUsers2(t *testing.T) {
	requestBody := models.IncomeCreateMultipleUsers{
		PayerIDs:     testUsersIDs[:3],
		ServiceIDs:   &[]openapitypes.UUID{serviceWDPObject.ID},
		PassInUseIDs: nil,
		Comment:      nil,
		Name:         &incomeNameMU3,
	}

	marshalledRequestBody, _ := json.Marshal(requestBody)

	responseBody := models.MessageOnlyResponse{}
	correctResponseBody := models.MessageOnlyResponse{
		Message: "Incomes were created successfully",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/incomes/multiple_users", bytes.NewBuffer(marshalledRequestBody))
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, correctResponseBody, responseBody)
}
func TestCreateIncomeDPMultipleUsers2Check(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      (*serviceWDPObject.DynamicPrices)[0].Price,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWDPObject.ID,
		Service:     &serviceWDPObject,
		Name:        &incomeNameMU3,
		IsPaid:      false,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	for _, user := range testUsers[:3] {
		correctResponseBody.PayerID = user.ID
		checkIncomeEquality(t, correctResponseBody, responseBody, false, false, testUsers[:3])
	}
}
func TestCreateIncomeDPMultipleUsers2CheckWrong(t *testing.T) {
	var responseBody []models.Income
	correctResponseBody := models.Income{
		IsActive:    true,
		Amount:      (*serviceWDPObject.DynamicPrices)[1].Price,
		Comment:     nil,
		UserID:      adminUserID,
		PassInUseID: nil,
		PassInUse:   nil,
		ServiceID:   &serviceWDPObject.ID,
		Service:     &serviceWDPObject,
		Name:        &incomeNameMU3,
		IsPaid:      false,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/incomes", nil)
	req.AddCookie(adminCookies[0])
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	for _, user := range testUsers[:3] {
		correctResponseBody.PayerID = user.ID
		checkIncomeEquality(t, correctResponseBody, responseBody, false, true, testUsers[:3])
	}
}
