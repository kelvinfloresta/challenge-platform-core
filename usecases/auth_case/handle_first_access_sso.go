package auth_case

import (
	"conformity-core/frameworks/logger"

	"github.com/sirupsen/logrus"
)

func (a AuthCase) handeFirstAccessSSO(login string) error {
	log := logger.Internal.WithFields(logrus.Fields{"login": login})
	workspace, err := a.getWorkspaceByLogin(login)
	if err != nil {
		return err
	}

	if workspace != "" {
		return SSORequired(workspace)
	}

	log.Info("login-not-found")
	return nil
}
