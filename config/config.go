package config

import "os"

type DatabaseConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
}

var (
	LoggerLevel         = os.Getenv("LOG_LEVEL")
	AppURL              = os.Getenv("APP_URL")
	APIHost             = os.Getenv("API_HOST")
	MailDomain          = os.Getenv("MAIL_DOMAIN")
	MailKey             = os.Getenv("MAIL_KEY")
	MailSender          = os.Getenv("MAIL_SENDER")
	BillingEmail        = os.Getenv("BILLING_MAIL")
	AuthSecret          = os.Getenv("AUTH_SECRET")
	ResetPasswordSecret = os.Getenv("RESET_PASSWORD_SECRET")
)

var Database = &DatabaseConfig{
	Host:         os.Getenv("DATABASE_HOST"),
	DatabaseName: os.Getenv("DATABASE_NAME"),
	User:         os.Getenv("DATABASE_USER"),
	Password:     os.Getenv("DATABASE_PASSWORD"),
	Port:         os.Getenv("DATABASE_PORT"),
}
