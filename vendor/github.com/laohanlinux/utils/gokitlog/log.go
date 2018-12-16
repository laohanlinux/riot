package gokitlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	logNameFormat = "2006-01-02_15:04"
	CallerNum     = 5
)

type LogFormat byte

const (
	JsonFormat LogFormat = 1 << iota
	FmtFormat
	NoFormat
)

func init() {
	tmpLog := log.NewJSONLogger(os.Stdout)
	tmpLog = log.With(tmpLog, "caller", log.DefaultCaller)
	tmpLog = log.With(tmpLog, "ts", log.DefaultTimestampUTC)
	tmpLog = level.NewFilter(tmpLog, level.AllowAll())
	tmpLog = log.NewSyncLogger(tmpLog)

	lg = &GoKitLogger{
		Logger:   tmpLog,
		ioWriter: nil,
		sync:     true,
	}
}

func NewGoKitLogger(opt LogOption, format ...LogFormat) (*GoKitLogger, error) {
	var (
		ioWriter *LogWriter
		tmpLog   log.Logger
		err      error
	)

	if len(format) == 0 {
		if ioWriter, err = NewLogWriter(opt); err != nil {
			return nil, err
		}
		tmpLog = log.NewLogfmtLogger(ioWriter)
	} else {
		switch format[0] {
		case JsonFormat:
			if ioWriter, err = NewLogWriter(opt); err != nil {
				return nil, err
			}
			tmpLog = log.NewJSONLogger(ioWriter)
		case FmtFormat:
			if ioWriter, err = NewLogWriter(opt); err != nil {
				return nil, err
			}
			tmpLog = log.NewLogfmtLogger(ioWriter)
		case NoFormat:
			ioWriter = &LogWriter{File: nil}
			tmpLog = log.NewNopLogger()
		default:
			panic("log format type is invalid")
		}
	}
	tmpLog = log.With(tmpLog, "caller", log.DefaultCaller)
	tmpLog = log.With(tmpLog, "ts", log.DefaultTimestampUTC)

	var gokitOpt level.Option
	switch strings.ToLower(opt.LogLevel) {
	case "info":
		gokitOpt = level.AllowInfo()
	case "debug":
		gokitOpt = level.AllowDebug()
	case "warn":
		gokitOpt = level.AllowWarn()
	case "error", "crit":
		gokitOpt = level.AllowError()
	default:
		panic(fmt.Sprintf("logLevel(%s) no in [info|debug|warn|error]", opt.LogLevel))
	}
	tmpLog = level.NewFilter(tmpLog, gokitOpt)

	if opt.Sync {
		tmpLog = log.NewSyncLogger(tmpLog)
	}
	return &GoKitLogger{Logger: tmpLog}, nil
}

func GlobalLog() *GoKitLogger {
	return lg
}

// it not thread saftly
func SetGlobalLog(opt LogOption) {
	Close()
	tmpLog, err := NewGoKitLogger(opt)
	if err != nil {
		panic(err)
	}
	lg = tmpLog
}

// SetGlobalLogWithLog set global logger with args logger
func SetGlobalLogWithLog(logger log.Logger, levelConf ...level.Option) {
	ioWriter := lg.ioWriter
	defer func() {
		if ioWriter != nil {
			ioWriter.Close()
		}
	}()
	// new logger has a new ioWriter
	lg.Logger = log.With(logger, "caller", log.Caller(CallerNum))
	if len(levelConf) > 0 {
		lg.Logger = level.NewFilter(lg.Logger, levelConf[0])
	} else {
		lg.Logger = level.NewFilter(lg.Logger, levelConf[0])
	}
}

var lg *GoKitLogger

type GoKitLogger struct {
	log.Logger
	// *levels.Levels
	ioWriter *LogWriter
	sync     bool
}

func (gklog *GoKitLogger) Close() error {
	if gklog.ioWriter != nil {
		return gklog.ioWriter.Close()
	}
	return nil
}

func Debug(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.DebugValue())
	logPrint(tmpLog, args)
}

func Debugf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.DebugValue())
	logPrintf(tmpLog, args)
}

func Info(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.InfoValue())
	logPrint(tmpLog, args)
}

func Infof(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.InfoValue())
	logPrintf(tmpLog, args)
}

func Warn(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.WarnValue())
	logPrint(tmpLog, args)
}

func Warnf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.WarnValue())
	logPrintf(tmpLog, args)
}

func Error(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
}

func Errorf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrintf(tmpLog, args)
}

func Crit(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
	os.Exit(1)
}

func Critf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
	os.Exit(1)
}

func Log(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Close() error {
	if lg.ioWriter != nil {
		return lg.ioWriter.Close()
	}
	return nil
}

func WrapLogLevel(levelsSets []string) []level.Option {
	var gokitOpts []level.Option
	for _, logLevel := range levelsSets {
		switch logLevel {
		case "info":
			gokitOpts = append(gokitOpts, level.AllowInfo())
		case "debug":
			gokitOpts = append(gokitOpts, level.AllowDebug())
		case "warn":
			gokitOpts = append(gokitOpts, level.AllowWarn())
		case "error":
			gokitOpts = append(gokitOpts, level.AllowError())
		}
	}
	return gokitOpts
}

type LogOption struct {
	// unit in minutes
	SegmentationThreshold int    `toml:"threshold"`
	LogDir                string `toml:"log_dir"`
	LogName               string `toml:"log_name"`
	LogLevel              string `toml:"log_level"`
	Sync                  bool   `toml:"sync"`
}

type LogWriter struct {
	oldTime               time.Time
	segmentationThreshold float64
	logDir                string
	logName               string
	*os.File
}

func NewLogWriter(opt LogOption) (*LogWriter, error) {
	logWriter := &LogWriter{
		oldTime:               time.Now(),
		segmentationThreshold: float64(opt.SegmentationThreshold),
		logDir:                opt.LogDir,
		logName:               opt.LogName,
	}

	fp, err := os.OpenFile(fmt.Sprintf("%s/%s.log", opt.LogDir, opt.LogName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	logWriter.File = fp
	return logWriter, nil
}

// TODO
// use bufio buffer
func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if lw.File == nil {
		return
	}
	if time.Since(lw.oldTime).Minutes() > lw.segmentationThreshold {
		if err = lw.renameLogFile(); err != nil {
			return -1, err
		}

		lw.File, err = os.OpenFile(fmt.Sprintf("%s/%s.log", lw.logDir, lw.logName),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return -1, err
		}

	}
	return lw.File.Write(p)
}

func (lw *LogWriter) Close() error {
	if lw.File != nil {
		return lw.File.Close()
	}
	return nil
}

func (lw *LogWriter) renameLogFile() (err error) {
	var (
		stat                     os.FileInfo
		srcFileName, dstFileName string
	)

	if lw.File == nil {
		return
	}

	if stat, err = lw.File.Stat(); err != nil {
		return err
	}
	srcFileName = fmt.Sprintf("%s/%s", lw.logDir, stat.Name())
	if err = lw.File.Close(); err != nil {
		return err
	}
	dstFileName = fmt.Sprintf("%s/%s_%s.log", lw.logDir, lw.logName,
		lw.oldTime.Format(logNameFormat))
	fmt.Println(dstFileName, srcFileName)
	os.Rename(srcFileName, dstFileName)
	lw.oldTime = time.Now()
	return nil
}

func logPrint(logger log.Logger, args []interface{}) {
	if args == nil || len(args) == 0 {
		logger.Log()
		return
	}
	if len(args) == 1 {
		logger.Log("msg", fmt.Sprintf("%v", args[0]))
		return
	}
	for i := 0; i < len(args); i++ {
		args[i] = fmt.Sprintf("%v", args[i])
	}
	logger.Log(args...)
}

func logPrintf(logger log.Logger, args []interface{}) {
	var (
		logFormat, msgContent string
	)
	if args == nil || len(args) == 0 {
		logger.Log("msg")
		return
	}
	if len(args) == 1 {
		logger.Log("msg", fmt.Sprintf("%s", args[0]))
		return
	}

	logFormat = fmt.Sprintf("%v", args[0])
	msgContent = fmt.Sprintf(logFormat, args[1:]...)
	logger.Log("msg", msgContent)
}
