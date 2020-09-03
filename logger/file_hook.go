package logger

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func newLfsHook(logFilePath, s string, formatter log.Formatter) log.Hook {
	filename := "%Y-%m-%d.log"
	if s != "" {
		filename = s + "." + filename
	}

	if logFilePath != "" {
		filename = path.Join(logFilePath, filename)
	}

	writer, _ := rotatelogs.New(
		filename,
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		// rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间
		rotatelogs.WithRotationTime(24*time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		// rotatelogs.WithMaxAge(time.Hour*24),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.DebugLevel: writer,
	}, formatter)

	return lfsHook
}
