package logger

import (
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

const transportProtocol = "tcp"

const (
	Ltimestamp = 1 << iota
	LJSON
)

const (
	connected = iota
	disconnected
)

type LoggerType struct {
	sync.Mutex
	logger        *logrus.Logger
	contextLogger *logrus.Entry
	conn          *net.Conn
	name          string
	address       string
	status        int32
}

func New(name, address string, flags int) (*LoggerType, error) {

	loggerType := &LoggerType{
		logger: &logrus.Logger{
			Out:       os.Stdout, // write log to stdout by default
			Formatter: makeFormatter(flags),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
		name:    name,
		address: address,
		status:  disconnected,
	}

	loggerType.contextLogger = logrus.NewEntry(loggerType.logger)
	loggerType.contextLogger = loggerType.contextLogger.WithField("module", name)

	if len(address) != 0 {
		if conn, err := net.Dial(transportProtocol, address); err == nil {
			loggerType.conn = &conn
			loggerType.logger.SetOutput(conn)
			loggerType.status = connected
		} else {
			return nil, err
		}
	}

	return loggerType, nil
}

func makeFormatter(flags int) (formatter logrus.Formatter) {

	disableTimestamp := (flags&Ltimestamp == 0)

	if (flags & LJSON) == 0 {
		formatter = &logrus.TextFormatter{
			DisableColors:          false,
			QuoteEmptyFields:       true,
			DisableTimestamp:       disableTimestamp,
			FullTimestamp:          true,
			TimestampFormat:        "15:04:05.000000",
			DisableLevelTruncation: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "xtime",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message"},
		}
	} else {
		formatter = &logrus.JSONFormatter{
			DisableTimestamp: disableTimestamp,
			TimestampFormat:  "15:04:05.000000",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "xtime",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message"},
		}
	}
	return
}

/*
GetName() - return module name
*/
func (lt *LoggerType) GetName() string {
	return lt.name
}

/*
SetLogLevel set level of logging. Severity log level can be one of: "debug", "info", "warn", "warning", "error", "fatal", "panic".
*/
func (lt *LoggerType) SetLogLevel(level string) (err error) {
	var lvl logrus.Level
	if lvl, err = logrus.ParseLevel(level); err != nil {
		return
	}

	lt.logger.SetLevel(lvl)
	return nil
}

func (lt *LoggerType) GetLogLevel() string {
	return lt.logger.GetLevel().String()
}

func (lt *LoggerType) checkConnection() bool {
	if lt.conn == nil || atomic.LoadInt32(&lt.status) == disconnected {
		return false
	}
	var buffer []byte
	lt.Lock()
	defer lt.Unlock()
	if _, err := (*lt.conn).Read(buffer); err != nil {
		(*lt.conn).Close()
		atomic.StoreInt32(&lt.status, disconnected)
		return false
	}
	return true
}

func (lt *LoggerType) Addfield(key string, value interface{}) {
	lt.contextLogger = lt.contextLogger.WithField(key, value)
}

func (lt *LoggerType) Addfields(fileds map[string]interface{}) {
	lt.contextLogger = lt.contextLogger.WithFields(fileds)
}

func (lt *LoggerType) reconnect() bool {
	if lt.conn == nil {
		return false
	}
	lt.Lock()
	defer lt.Unlock()
	if atomic.LoadInt32(&lt.status) == connected {
		return true
	}
	(*lt.conn).Close()
	if conn, err := net.Dial(transportProtocol, lt.address); err == nil {
		lt.conn = &conn
		lt.logger.SetOutput(conn)
		atomic.StoreInt32(&lt.status, connected)
		return true
	}
	lt.logger.SetOutput(os.Stdout)
	return false
}

func (lt *LoggerType) checkAndReconnect() {
	if lt.conn != nil && lt.checkConnection() == false {
		lt.reconnect()
	}
}

func (lt *LoggerType) Debug(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Debug(v)
	}
}

func (lt *LoggerType) Info(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.InfoLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Info(v)
	}
}

func (lt *LoggerType) Warning(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.WarnLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Warn(v)
	}
}

func (lt *LoggerType) Error(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Error(v)
	}
}

func (lt *LoggerType) Fatal(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.FatalLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Fatal(v)
	}
}

func (lt *LoggerType) Panic(v ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.PanicLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Panic(v)
	}
}

func (lt *LoggerType) Debugf(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Debugf(format, args...)
	}
}

func (lt *LoggerType) Infof(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.InfoLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Infof(format, args...)
	}
}

func (lt *LoggerType) Warningf(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.WarnLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Warnf(format, args...)
	}
}

func (lt *LoggerType) Errorf(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Errorf(format, args...)
	}
}

func (lt *LoggerType) Fatalf(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.FatalLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Fatalf(format, args...)
	}
}

func (lt *LoggerType) Panicf(format string, args ...interface{}) {
	if lt.contextLogger.Logger.IsLevelEnabled(logrus.PanicLevel) {
		lt.checkAndReconnect()
		lt.contextLogger.Panicf(format, args...)
	}
}

/*
InitLogger - universal function for creation logger. Logger can write output to stdout or to remote host:port
To send logs to remote host use address parameter, for example "192.168.0.1:5001".
To send logs to stdout address parameter must be empty string.
Name parameter will be add to every log message like module="name"
Severity log level can be one of: "debug", "info", "warn", "warning", "error", "fatal", "panic".
*/
func InitLogger(name, severity, address string, flags int) (lg *LoggerType, err error) {
	if lg, err = New(name, address, flags); err != nil {
		return nil, err
	}
	if len(severity) == 0 {
		severity = "debug"
	}
	if err = lg.SetLogLevel(severity); err != nil {
		log.Println(err)
	}
	return
}
