package user_case

import (
	coreContext "conformity-core/context"
	core_errors "conformity-core/errors"
	"conformity-core/gateways/user_gateway"
	"conformity-core/utils"

	"github.com/sirupsen/logrus"
)

type ChangePasswordInput struct {
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
}

func (u *UserCase) ChangePassword(ctx *coreContext.CoreCtx, input *ChangePasswordInput) (bool, error) {
	if input.Password != input.PasswordConfirmation {
		return false, core_errors.PasswordNotMatch
	}

	encryptedPassword, err := utils.PasswordGen(input.Password)

	if err != nil {
		return false, err
	}

	changed, err := u.gateway.ChangePassword(user_gateway.ChangePasswordInput{
		UserID:   ctx.Session.UserID,
		Password: encryptedPassword,
	})

	if err != nil {
		return false, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"changed": changed,
	}).Info("password-changed")

	return changed, nil
}
