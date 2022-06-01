package notification_gateway

import (
	"conformity-core/context"
	"time"
)

type Attachment struct {
	Name string
	File []byte
}

type SendEmailInput struct {
	To         []string
	Subject    string
	TemplateID string
	Schedule   *time.Time
	Variables  map[string]string
	Attachment *Attachment
}

type SendEmail func(ctx *context.CoreCtx, input SendEmailInput) error
