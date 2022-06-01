package user_case

import (
	"conformity-core/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func generateResetPasswordToken(userId string, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now()
		now = &t
	}

	claims := jwt.MapClaims{
		"id":  userId,
		"exp": now.Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.ResetPasswordSecret))

	return t, err
}
