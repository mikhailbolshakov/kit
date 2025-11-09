package http

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"net/http"
	"testing"
)

type decodeRequestTestSuite struct {
	kit.Suite
}

func (s *decodeRequestTestSuite) SetupSuite() {
	s.Suite.Init(logf)
}

func (s *decodeRequestTestSuite) SetupTest() {
}

func TestDecodeRequestSuite(t *testing.T) {
	suite.Run(t, new(decodeRequestTestSuite))
}

type testRequest struct {
	Name        string  `json:"name" validate:"required" san:"trim"`
	Description *string `json:"description" san:"trim,max=50"`
	Email       string  `json:"email" validate:"required,email" san:"trim,lower"`
	PhoneNumber string  `json:"phone_number" validate:"required,e164"`
	Title       string  `json:"title" san:"trim,max=20"`
	Department  *string `json:"department" san:"trim,def=General"`
	Code        string  `json:"code" validate:"required" san:"trim,upper"`
	Age         int     `json:"age" validate:"required,gte=0,lte=150"`
	Score       float64 `json:"score" validate:"gte=0,lte=100"`
	Count       *int    `json:"count" validate:"omitempty,gte=0"`
	Active      bool    `json:"active"`
}

func (s *decodeRequestTestSuite) Test_ValidRequest() {
	jsonStr := `{
		"name": "  John Doe  ",
		"description": "   Test Description   ",
		"email": "   TEST@EXAMPLE.COM   ",
		"phone_number": "+12125552368",
		"title": "   Senior Developer    ",
		"department": null,
		"code": "abc123",
		"age": 30,
		"score": 85.5,
		"count": 5,
		"active": true
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	decoded, err := DecodeRequest[testRequest](s.Ctx, req)

	s.NoError(err)
	s.NotNil(decoded)

	s.Equal("John Doe", decoded.Name)
	s.Equal("Test Description", *decoded.Description)
	s.Equal("test@example.com", decoded.Email)
	s.Equal("Senior Developer", decoded.Title)
	s.Equal("General", *decoded.Department)
	s.Equal("ABC123", decoded.Code)
	s.Equal(30, decoded.Age)
	s.Equal(85.5, decoded.Score)
	s.Equal(5, *decoded.Count)
	s.True(decoded.Active)
}

func (s *decodeRequestTestSuite) Test_MissingRequiredFields() {
	jsonStr := `{
		"description": "Test",
		"title": "Developer",
		"score": 85.5
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}

func (s *decodeRequestTestSuite) Test_InvalidEmail() {
	jsonStr := `{
		"name": "John Doe",
		"email": "invalid-email",
		"phone_number": "+12125552368",
		"code": "abc123",
		"age": 30
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}

func (s *decodeRequestTestSuite) Test_InvalidPhoneNumber() {
	jsonStr := `{
		"name": "John Doe",
		"email": "john@example.com",
		"phone_number": "invalid-phone",
		"code": "abc123",
		"age": 30
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}

func (s *decodeRequestTestSuite) Test_NumericValidation() {
	jsonStr := `{
		"name": "John Doe",
		"email": "john@example.com",
		"phone_number": "+12125552368",
		"code": "abc123",
		"age": 200,
		"score": 150,
		"count": -1
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}

func (s *decodeRequestTestSuite) Test_MaxLengthValidation() {
	description := "This is a very long description that exceeds the maximum length of fifty characters"
	jsonStr := `{
		"name": "John Doe",
		"description": "` + description + `",
		"email": "john@example.com",
		"phone_number": "+12125552368",
		"code": "abc123",
		"age": 30
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	decoded, err := DecodeRequest[testRequest](s.Ctx, req)

	s.NoError(err)
	s.NotNil(decoded)
	s.Equal(50, len(*decoded.Description))
}

func (s *decodeRequestTestSuite) Test_InvalidJSON() {
	jsonStr := `{
		"name": "John Doe",
		"email": "john@example.com",
		"phone_number": "+12125552368",
		"code": "abc123",
		"age": "invalid", 
	}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpDecodeRequest)
}

func (s *decodeRequestTestSuite) Test_EmptyRequest() {
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(""))
	_, err := DecodeRequest[testRequest](s.Ctx, req)
	s.AssertAppErr(err, ErrCodeHttpDecodeRequest)
}

type testUUIDRequest struct {
	ID         string  `json:"id" validate:"required,uuid" san:"trim"`
	ParentID   *string `json:"parent_id" validate:"omitempty,uuid" san:"trim"`
	OptionalID string  `json:"optional_id" validate:"omitempty,uuid" san:"trim"`
}

// Add these test cases to the suite
func (s *decodeRequestTestSuite) Test_UUIDWithTrailingSpaces() {
	jsonStr := `{
        "name": "John Doe",
        "email": "john@example.com",
        "phone_number": "+12125552368",
        "code": "abc123",
        "age": 30,
        "id": "  550e8400-e29b-41d4-a716-446655440000  "
    }`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	decoded, err := DecodeRequest[testUUIDRequest](s.Ctx, req)

	s.NoError(err)
	s.NotNil(decoded)
	s.Equal("550e8400-e29b-41d4-a716-446655440000", decoded.ID)
}

func (s *decodeRequestTestSuite) Test_EmptyUUIDWithOmitempty() {
	jsonStr := `{
        "name": "John Doe",
        "email": "john@example.com",
        "phone_number": "+12125552368",
        "code": "abc123",
        "age": 30,
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "optional_id": "",
        "parent_id": null
    }`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	decoded, err := DecodeRequest[testUUIDRequest](s.Ctx, req)

	s.NoError(err)
	s.NotNil(decoded)
	s.Empty(decoded.OptionalID)
	s.Nil(decoded.ParentID)
}

func (s *decodeRequestTestSuite) Test_InvalidUUID() {
	jsonStr := `{
        "name": "John Doe",
        "email": "john@example.com",
        "phone_number": "+12125552368",
        "code": "abc123",
        "age": 30,
        "id": "not-a-uuid"
    }`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testUUIDRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}

func (s *decodeRequestTestSuite) Test_InvalidUUIDFormat() {
	jsonStr := `{
        "name": "John Doe",
        "email": "john@example.com",
        "phone_number": "+12125552368",
        "code": "abc123",
        "age": 30,
        "id": "550e8400-e29b-41d4-a716-44665544000" 
    }`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	_, err := DecodeRequest[testUUIDRequest](s.Ctx, req)

	s.AssertAppErr(err, ErrCodeHttpValidationRequest)
}
