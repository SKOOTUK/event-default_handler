package utils

import (
	"log"
	"os"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
)

// FailOnError raise to Sentry and fail on error
func FailOnError(err error, msg string) {
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("reported to Sentry: %s", err)
		log.Fatalf("%s: %s", msg, err)
	}
}

// SetupSentry run initial Sentry configuration
func SetupSentry() {
	sampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_SAMPLE_RATE"), 64)
	err = sentry.Init(sentry.ClientOptions{
		Dsn:        os.Getenv("SENTRY_DSN"),
		SampleRate: sampleRate,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
