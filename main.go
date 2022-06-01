package main

import (
	database "conformity-core/frameworks/database/gorm"
	app "conformity-core/frameworks/http/fiber"
	"conformity-core/frameworks/scheduler"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	time.Local = time.UTC
	err := sentry.Init(sentry.ClientOptions{
		TracesSampleRate: 0.2,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	database.CreateDatabase()
	database.DB_Production.Connect()
	if migrationError := database.Migrate(database.DB_Production.Conn); migrationError != nil {
		sentry.CaptureException(migrationError)
		log.Fatal(migrationError)
	}

	err = scheduler.Start()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	app := app.CreateApp()
	log.Fatal(app.Listen(":3001"))
}
