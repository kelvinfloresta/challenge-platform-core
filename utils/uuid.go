package utils

import "github.com/google/uuid"

func UUID() (string, error) {
	u, err := uuid.NewRandom()
	if err == nil {
		return u.String(), nil
	}

	return "", err
}
