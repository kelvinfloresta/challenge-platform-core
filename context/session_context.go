package context

import (
	"conformity-core/enums/user_company_enum"
	"context"
)

type UserSession struct {
	UserID         string
	Name           string
	Login          string
	DepartmentID   string
	CompanyID      string
	OriginalUserID string
	Role           user_company_enum.Role
}

const SessionKey = "session"

func getSession(ctx context.Context) *UserSession {
	session, ok := ctx.Value(SessionKey).(*UserSession)
	if !ok {
		return &UserSession{}
	}

	return session
}
