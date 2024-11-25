package retry

import (
	"errors"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
)

type Func func(retryCount int) error

type ShouldStopFunc func() bool

type OnFailedAttemptFunc func(error)

type DelayFunc func(int) time.Duration

func Do(retryFunc Func, stopFunc ShouldStopFunc, delayFunc DelayFunc, attemptFunc OnFailedAttemptFunc) error {
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
		var retryErr *agora.RetryErr
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
