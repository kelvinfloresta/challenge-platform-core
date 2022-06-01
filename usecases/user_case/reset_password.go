package user_case

import (
	"conformity-core/context"
	core_errors "conformity-core/errors"
	"conformity-core/frameworks/logger"
	"conformity-core/usecases/notification_case"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ResetPasswordInput struct {
	Email string `json:"email" validate:"required"`
}

func (u *UserCase) ResetPassword(ctx *context.CoreCtx, input *ResetPasswordInput) error {
	user, err := u.GetByEmail(input.Email)

	if err != nil {
		return fmt.Errorf("get by email: %v", err)
	}

	if user == nil {
		return core_errors.UserNotFound
	}

	token, err := generateResetPasswordToken(user.ID, nil)

	if err != nil {
		return fmt.Errorf("generate token: %v", err)
	}

	err = u.notificationCase.ResetPassword(ctx, notification_case.ResetPasswordInput{
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	})

	if err != nil {
		return fmt.Errorf("notify reset password: %v", err)
	}

	logger.Internal.WithFields(logrus.Fields{
		"notification": "reset-password",
		"email":        input.Email,
	}).Info("task-done")

	return nil
}
