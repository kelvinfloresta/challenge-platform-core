package department_controller_test

import (
	"bytes"
	"conformity-core/fixtures"
	app "conformity-core/frameworks/http/fiber"
	"conformity-core/usecases/department_case"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func postDepartment(input department_case.CreateDepartmentCaseInput, app *fiber.App) string {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
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

	data := department_case.CreateDepartmentCaseInput{
		Name:      "",
		CompanyID: "any",
	}

	// When
	response := postDepartment(data, app)
	expectedResult := `{"errors":[{"FailedField":"Name","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}

func TestCompanyIdRequired(t *testing.T) {
	t.SkipNow()
	// Given
	app := app.CreateApp()

	data := department_case.CreateDepartmentCaseInput{
		Name: "",
	}

	// When
	response := postDepartment(data, app)
	expectedResult := `{"errors":[{"FailedField":"CompanyID","Tag":"required","Value":""}]}`

	// Verify
	if expectedResult != response {
		t.Errorf("Result %s; Want %s", response, expectedResult)
	}
}
