package core

import (
	"errors"
	"time"
)

type RetryFunc func(retryCount int) error

type RetryShouldStopFunc func() bool

type OnFailedAttemptFunc func(error)

type DelayFunc func(int) time.Duration

func RetryDo(retryFunc RetryFunc, stopFunc RetryShouldStopFunc, delayFunc DelayFunc, attemptFunc OnFailedAttemptFunc) error {
	var (
		retryCount int
		err        error
	)

	stopRetry := stopFunc()
	for !stopRetry {
		err = retryFunc(retryCount)
		if err == nil {
			return nil
		}
		var retryErr *RetryErr
		if errors.As(err, &retryErr) {
			if !retryErr.NeedRetry() {
				return retryErr
			}
		}

		if attemptFunc != nil {
			attemptFunc(err)
		}
		stopRetry = stopFunc()
		if !stopRetry && delayFunc != nil {
			time.Sleep(delayFunc(retryCount))
		}
		retryCount++
	}
	return err
}
