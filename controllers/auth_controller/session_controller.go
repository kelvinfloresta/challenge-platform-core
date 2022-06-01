package auth_controller

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Session(c *fiber.Ctx) error {
	token, ok := c.Context().UserValue("user").(*jwt.Token)
	if !ok || !token.Valid {
		return errors.New("INVALID_TOKEN")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("INVALID_CLAIMS")
	}

	if err := claims.Valid(); err != nil {
		return errors.New("INVALID_CLAIMS")
	}

	c.Context().SetUserValue(string(context.SessionKey), &context.UserSession{
		UserID:       fmt.Sprint(claims["id"]),
		Name:         fmt.Sprint(claims["name"]),
		Login:        fmt.Sprint(claims["login"]),
		DepartmentID: fmt.Sprint(claims["department_id"]),
		CompanyID:    fmt.Sprint(claims["company_id"]),
		Role:         user_company_enum.Role(fmt.Sprint(claims["role"])),
	})

	return c.Next()
}
