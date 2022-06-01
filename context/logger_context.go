package context

import (
	"conformity-core/frameworks/logger"
	"context"

	"github.com/sirupsen/logrus"
)

const requestIdKey = "requestid"

var log = logger.New()

func getRequestId(ctx context.Context) string {
	id, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		return ""
	}

	return id
}

func getLogger(ctx context.Context, session *UserSession) *logrus.Entry {
	if session == nil {
		return log.WithField(requestIdKey, getRequestId(ctx))
	}
	entry := log.WithFields(logrus.Fields{
		requestIdKey:   getRequestId(ctx),
		"userId":       session.UserID,
		"departmentId": session.DepartmentID,
	})

	if session.Role.IsBackoffice() {
		entry.WithFields(logrus.Fields{
			"originalUserId": session.OriginalUserID,
		})
	}

	return entry
}
