package logger

import (
	"bufio"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(logPath string, state bool) {

	hit, err := isFileExist(logPath)
	if !hit {
		err = os.MkdirAll(logPath, os.ModePerm)
	}
	if err != nil {
		panic(err)
	}

	logMap := make(map[string]logrus.Level)
	logMap["debug"] = logrus.DebugLevel
	logMap["info"] = logrus.InfoLevel
	logMap["warn"] = logrus.WarnLevel
	logMap["error"] = logrus.ErrorLevel
	logMap["fatal"] = logrus.FatalLevel
	logMap["panic"] = logrus.PanicLevel

	logWriteMap := lfshook.WriterMap{}
	logFormatter := &logrus.JSONFormatter{}
	maxAge := 2 * 24 * time.Hour
	rotationTime := 24 * time.Hour

	for key, level := range logMap {
		baseLogPaht := path.Join(logPath, key)
		writer, err := rotatelogs.New(
			baseLogPaht+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
			rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
		)
		if err != nil {
			logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
		}
		logWriteMap[level] = writer
	}

	lfHook := lfshook.NewHook(logWriteMap, logFormatter)
	logrus.AddHook(lfHook)

	// 不打印日志到控制台了
	SetConsole(state)
	// end
}

// 判断文件文件夹是否存在
func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func SetConsole(state bool) {
	if !state {
		// 不打印日志到控制台了
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logrus.Errorf("Err: %+v", errors.WithStack(err))
		}
		writer := bufio.NewWriter(src)
		logrus.SetOutput(writer)
	}
}

func GetLogger() *logrus.Logger {

	return logrus.StandardLogger()
}

func Printf(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Printf(format, args...)
}

func Debug(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	//Info(format, args...)
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Infof(format, args...)
}

func Fatal(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Panicf(format, args...)
	//Panic(format, args...)
}

func Info(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Errorf(format, args...)
}

func Panic(format string, args ...interface{}) {
	if args == nil {
		args = append(args, "")
		format = format + "%s"
	}
	// Calls panic() after logging
	logrus.WithFields(logrus.Fields{
		//"time":       time.Now(),
	}).Panicf(format, args...)
}
