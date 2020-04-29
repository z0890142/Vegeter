package Comman

import (
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var Log *logrus.Entry

func LogInit(packageName string, projectName string, level logrus.Level) (*logrus.Entry, string) {
	msgId := uuid.NewV4()

	log = logrus.New()
	log.SetLevel(level)       //設定log level
	log.SetReportCaller(true) //開啟列印log位置資訊 可能會導致效能下降
	log.SetFormatter(&logrus.JSONFormatter{
		// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 	filePathArray := strings.Split(f.File, "/")
		// 	fileInfo := filePathArray[len(filePathArray)-1] + " : " + strconv.Itoa(f.Line)
		// 	return f.Function, fileInfo
		// },
		TimestampFormat: "2006-01-02T15:04:05.000000",
		PrettyPrint:     true,
	})
	_dir := path.Join("logs", projectName, packageName)
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
	}
	if exist {
		fmt.Printf("has dir![%v]\n", _dir)
	} else {
		fmt.Printf("no dir![%v]\n", _dir)
		// Mkdir 建立的資料夾前面的父資料夾必須全都存在，只要有一個不存在會返回err
		// MkdirAll 建立的資料夾前面的父資料夾如果不存在，會一併建立
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}

	ConfigLocalFilesystemLogger(_dir, packageName, time.Hour*24*30, time.Minute*30)
	Log = log.WithFields(logrus.Fields{
		"PackageName": packageName,
		"LogSN":       fmt.Sprintf("%s", msgId),
	})
	return Log, fmt.Sprintf("%s", msgId)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+"_%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(baseLogPaht),
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存時間
		rotatelogs.WithRotationTime(rotationTime), // log切割的時間間格
	)
	if err != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 不同的級別輸出的位置
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 	filePathArray := strings.Split(f.File, "/")
		// 	fileInfo := filePathArray[len(filePathArray)-1] + " : " + strconv.Itoa(f.Line)
		// 	return f.Function, fileInfo
		// },
		TimestampFormat: "2006-01-02T15:04:05.000000",
		PrettyPrint:     false,
	})

	log.AddHook(lfHook)
}

//判斷資料夾是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
