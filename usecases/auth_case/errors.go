package auth_case

import (
	core_errors "conformity-core/errors"
	"fmt"
)

var (
	ErrPasswordRequired core_errors.Unauthorized = "PASSWORD_REQUIRED"
	ErrSuspendedUser    core_errors.Forbidden    = "SUSPENDED_USER"
	ErrWrongPassword    core_errors.Unauthorized = "WRONG_PASSWORD"
)

func SSORequired(workspace string) error {
	return core_errors.Unauthorized(fmt.Sprintf("SSO_REQUIRED__%s", workspace))
}
