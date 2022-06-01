package company_controller_test

import (
	"bytes"
	"conformity-core/fixtures"
	app "conformity-core/frameworks/http/fiber"
	"conformity-core/usecases/company_case"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func postCompany(input company_case.CreateCompanyCaseInput, app *fiber.App) string {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/companies", bytes.NewReader(body))
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

	data := company_case.CreateCompanyCaseInput{
		Name:     "",
		Document: "123456789",
	}

	// When
	response := postCompany(data, app)
	expectedResult := `{"errors":[{"FailedField":"Name","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestDocumentRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := company_case.CreateCompanyCaseInput{
		Name:     "Name",
		Document: "",
	}

	// When
	response := postCompany(data, app)
	expectedResult := `{"errors":[{"FailedField":"Document","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}
