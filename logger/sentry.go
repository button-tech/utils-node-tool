package logger

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"time"
)

func InitSentry(dsn string) error {
	if len(dsn) == 0 {
		return errors.New("Set DSN!")
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func SendError(err error) {
	sentry.CaptureException(err)
}

func SendFatal(err error) {
	sentry.CaptureException(err)
	if sentry.Flush(time.Second * 5) {
		Info("All sentry queued events delivered!")
	} else {
		Info("Sentry flush timeout reached")
	}
}

func SendMessage(msg string) {
	sentry.CaptureMessage(msg)
}
