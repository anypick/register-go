package baselog

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"register-go/infra"
	"register-go/infra/base/log/config"
	"register-go/infra/utils/common"
	"strings"
	"time"
)

type LogrusStarter struct {
	infra.BaseStarter
}

func (l *LogrusStarter) Init(ctx infra.StarterContext) {
	initLogrus(ctx.Yaml().LogConfig)
}

// 初始化Logrus
func initLogrus(config config.LogConfig) {
	var (
		level          logrus.Level                   // 日志级别
		logPath        = config.FilePath              // 日志存放路径
		fileName       = config.LogFileName           // 日志名称
		logPathAndName = path.Join(logPath, fileName) // 日志路径+名称
		out            *rotatelogs.RotateLogs         // 日志输出
		err            error
	)
	formatter := logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: common.TimeFormat,
	}
	logrus.SetFormatter(&formatter)
	switch config.Level {
	case "trace":
		level = logrus.TraceLevel
		break
	case "debug":
		level = logrus.DebugLevel
		break
	case "info":
		level = logrus.InfoLevel
		break
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
		break
	default:
		level = logrus.DebugLevel
		break
	}
	logrus.SetLevel(level)

	if level == logrus.TraceLevel {
		logrus.SetOutput(os.Stdout)
		return
	}
	// 创建日志目录，如果存在则忽略
	os.MkdirAll(logPath, os.ModePerm)
	// 定义日志切割
	out, err = rotatelogs.New(
		strings.TrimSuffix(logPathAndName, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(logPathAndName), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)),
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)),
	)
	if err != nil {
		logrus.Error("new rotatelogs error ", err)
		return
	}
	logrus.SetOutput(out)
}
