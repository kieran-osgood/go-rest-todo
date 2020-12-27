package error

import "go.uber.org/zap"

func CleanUpAndHandleError(cleanUp func() error, logger *zap.SugaredLogger) {
	if err := cleanUp(); err != nil {
		logger.Error(err)
	}
}

func CleanUpAndHandleErrorDefaultLogger(cleanUp func() error, logger *zap.Logger) {
	if err := cleanUp(); err != nil {
		logger.Error("Clean up failed")
	}
}
