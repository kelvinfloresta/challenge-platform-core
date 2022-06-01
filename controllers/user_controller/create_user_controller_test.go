package user_controller_test

import (
	"bytes"
	"conformity-core/fixtures"
	app "conformity-core/frameworks/http/fiber"
	"conformity-core/usecases/user_case"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func postUser(input user_case.CreateInput, app *fiber.App) string {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Add("Authorization", fixtures.TestToken)
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	response := buf.String()
	return response
}

func TestNameRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "",
		Password:    "Password",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Name","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestPasswordRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Password","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestEmailRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "Password",
		Email:       "",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Email","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestDocumentRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "Password",
		Email:       "me@email.com",
		Document:    "",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Document","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestPhoneRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "Password",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Phone","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestBirthDateRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "Password",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"BirthDate","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestJobPositionRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "Password",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"JobPosition","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestPasswordMin8Length(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := user_case.CreateInput{
		Name:        "Name",
		Password:    "1234567",
		Email:       "me@email.com",
		Document:    "Document",
		Phone:       "Phone",
		BirthDate:   "BirthDate",
		JobPosition: "JobPosition",
	}

	// When
	response := postUser(data, app)
	expectedResult := `{"errors":[{"FailedField":"Password","Tag":"min","Value":"8"}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}
