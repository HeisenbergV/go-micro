package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})

	Errorf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Debugf(f string, args ...interface{})
	WithFieldsInfo(fields map[string]interface{}, data string)
	WithFieldsError(fields map[string]interface{}, data string)
	WithFieldsDebug(fields map[string]interface{}, data string)
}

type loggerWrapper struct {
	lw *logrus.Logger
}

func (logger *loggerWrapper) WithFieldsInfo(fields map[string]interface{}, data string) {
	logger.lw.WithFields(logrus.Fields(fields)).Info(data)
}

func (logger *loggerWrapper) WithFieldsError(fields map[string]interface{}, data string) {
	logger.lw.WithFields(logrus.Fields(fields)).Error(data)
}

func (logger *loggerWrapper) WithFieldsDebug(fields map[string]interface{}, data string) {
	logger.lw.WithFields(logrus.Fields(fields)).Debug(data)
}

func (logger *loggerWrapper) Error(args ...interface{}) {
	logger.lw.Error(args...)
}

func (logger *loggerWrapper) Info(args ...interface{}) {
	logger.lw.Info(args...)
}
func (logger *loggerWrapper) Debug(args ...interface{}) {
	logger.lw.Debug(args...)
}

func (logger *loggerWrapper) Errorf(f string, args ...interface{}) {
	logger.lw.Errorf(f, args...)
}
func (logger *loggerWrapper) Infof(f string, args ...interface{}) {
	logger.lw.Infof(f, args...)
}
func (logger *loggerWrapper) Debugf(f string, args ...interface{}) {
	logger.lw.Debugf(f, args...)
}
