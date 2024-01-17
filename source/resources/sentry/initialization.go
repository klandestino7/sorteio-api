package sentry

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func SentryInitialization() {
	sentryConnection := os.Getenv("SENTRY_CONNECTION")
	err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryConnection,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
}
