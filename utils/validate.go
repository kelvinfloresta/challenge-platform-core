package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateStruct(s interface{}, c *fiber.Ctx) fiber.Map {
	err := validator.New().Struct(s)

	if err == nil {
		return nil
	}

	var errors []*ErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		element := ErrorResponse{
			FailedField: err.Field(),
			Tag:         err.Tag(),
			Value:       err.Param(),
		}
		errors = append(errors, &element)
	}

	c.Status(400)
	return fiber.Map{"errors": errors}
}

func parseStruct(data interface{}, c *fiber.Ctx) fiber.Map {
	parseError := c.BodyParser(data)

	if parseError == nil {
		return nil
	}

	c.Status(400)
	errors := []error{parseError}
	return fiber.Map{"errors": errors}

}

func ValidateBody(data interface{}, c *fiber.Ctx) fiber.Map {
	if parseError := parseStruct(data, c); parseError != nil {
		return parseError
	}

	if validateError := validateStruct(data, c); validateError != nil {
		return validateError
	}

	return nil
}
