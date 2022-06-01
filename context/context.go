package context

import (
	"conformity-core/enums/user_company_enum"
	"context"

	"github.com/sirupsen/logrus"
)

type CoreCtx struct {
	Logger  *logrus.Entry
	Context context.Context
	Session *UserSession
}

func New(ctx context.Context) *CoreCtx {
	session := getSession(ctx)
	return &CoreCtx{
		Context: ctx,
		Logger:  getLogger(ctx, session),
		Session: session,
	}
}

var Internal *CoreCtx

func init() {
	Internal = New(context.Background())
	Internal.Session = &UserSession{
		UserID:       "internal",
		Name:         "internal",
		Login:        "internal",
		DepartmentID: "internal",
		CompanyID:    "",
		Role:         user_company_enum.RoleBackoffice,
	}
}
