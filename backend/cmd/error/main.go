package error

import "go.uber.org/zap"

// CleanUpAndHandleError wraps a clean up function to be deferred, and logs any returned error
func CleanUpAndHandleError(cleanUp func() error, logger *zap.SugaredLogger) {
	if err := cleanUp(); err != nil {
		logger.Error(err)
	}
}
