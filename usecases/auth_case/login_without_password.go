package auth_case

import (
	"conformity-core/enums/user_company_enum"
	"conformity-core/frameworks/logger"

	"github.com/sirupsen/logrus"
)

func (a *AuthCase) LoginWithoutPassword(login string) (string, error) {
	user, err := a.departmentCase.GetUserByLogin(login)
	log := logger.Internal.WithFields(logrus.Fields{"login": login})

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", a.handeFirstAccessSSO(login)
	}

	if user.Status == user_company_enum.Suspended {
		log.Info("user-suspended")
		return "", ErrSuspendedUser
	}

	if user.SSOEnabled {
		return "", SSORequired(user.Workspace)
	}

	if user.RequirePassword {
		return "", ErrPasswordRequired
	}

	user.Role = user_company_enum.RoleUser
	return generateToken(user, nil)
}
