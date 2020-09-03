package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLog(service string) Logger {
	log := logrus.New()
	// if conf.Config.RunMode == "dev" {
	// 	log.SetLevel(logrus.DebugLevel)
	// 	log.SetOutput(os.Stdout)
	// } else {
	// 	log.SetLevel(logrus.InfoLevel)
	// }
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000", DisableColors: true})

	log.AddHook(newLfsHook(".", service, log.Formatter))
	wrapper := &loggerWrapper{log}
	return wrapper
}
