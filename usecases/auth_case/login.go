package auth_case

import (
	"conformity-core/enums/user_company_enum"
	"conformity-core/frameworks/logger"
	"conformity-core/utils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthCase) Login(input LoginInput) (string, error) {
	user, err := a.departmentCase.GetUserByLogin(input.Login)
	log := logger.Internal.WithFields(logrus.Fields{"login": input.Login})

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", a.handeFirstAccessSSO(input.Login)
	}

	if user.Status == user_company_enum.Suspended {
		log.Info("user-suspended")
		return "", ErrSuspendedUser
	}

	if user.SSOEnabled {
		return "", SSORequired(user.Workspace)
	}

	if user.Password == "" {
		log.Info("empty-password")
		return "", ErrWrongPassword
	}

	err = utils.ComparePassword(input.Password, user.Password)

	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Info("wrong-password")
		return "", ErrWrongPassword
	}

	if err != nil {
		return "", err
	}

	return generateToken(user, nil)
}
