package utils

import (
	"log"
	"os"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
)

// SetupSentry run initial Sentry configuration
func SetupSentry() {
	sampleRate, _ := strconv.ParseFloat(os.Getenv("SENTRY_SAMPLE_RATE"), 64)
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:        os.Getenv("SENTRY_DSN"),
		SampleRate: sampleRate,
	}); err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
