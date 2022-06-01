package auth_case

import (
	"conformity-core/config"
	"conformity-core/usecases/department_case"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func generateToken(user *department_case.UserDepartment, originalUser *string) (string, error) {
	claims := jwt.MapClaims{
		"id":            user.UserID,
		"name":          user.Name,
		"login":         user.Login,
		"role":          user.Role,
		"department_id": user.DepartmentID,
		"company_id":    user.CompanyID,
		"exp":           time.Now().Add(time.Hour * 72).Unix(),
	}

	if originalUser != nil {
		claims["original_user"] = originalUser
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.AuthSecret))

	return t, err
}
