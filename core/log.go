package core

import (
	"context"
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarning
	LogErr
	LogUnknown
)

type Logger interface {
	Debug(ctx context.Context, module string, v ...interface{})
	Debugf(ctx context.Context, module string, format string, v ...interface{})
	Error(ctx context.Context, module string, v ...interface{})
	Errorf(ctx context.Context, module string, format string, v ...interface{})
	Info(ctx context.Context, module string, v ...interface{})
	Infof(ctx context.Context, module string, format string, v ...interface{})
	Warn(ctx context.Context, module string, v ...interface{})
	Warnf(ctx context.Context, module string, format string, v ...interface{})

	Level() LogLevel
	SetLevel(level LogLevel)
}

var (
	_             Logger = (*discardLogger)(nil)
	DiscardLogger        = &discardLogger{}
)

type discardLogger struct{}

func (d *discardLogger) Debug(ctx context.Context, module string, v ...interface{}) {}

func (d *discardLogger) Debugf(ctx context.Context, module string, format string, v ...interface{}) {}

func (d *discardLogger) Error(ctx context.Context, module string, v ...interface{}) {}

func (d *discardLogger) Errorf(ctx context.Context, module string, format string, v ...interface{}) {}

func (d *discardLogger) Info(ctx context.Context, module string, v ...interface{}) {}

func (d *discardLogger) Infof(ctx context.Context, module string, format string, v ...interface{}) {}

func (d *discardLogger) Warn(ctx context.Context, module string, v ...interface{}) {}

func (d *discardLogger) Warnf(ctx context.Context, module string, format string, v ...interface{}) {}

func (d *discardLogger) Level() LogLevel {
	return LogUnknown
}

func (d *discardLogger) SetLevel(level LogLevel) {}

type sampleLogger struct {
	DEBUG *log.Logger
	ERROR *log.Logger
	INFO  *log.Logger
	WARN  *log.Logger

	level LogLevel
}

const (
	defaultLogPrefix = "agora-rest-client-go"
	defaultLogFlag   = log.Lshortfile | log.LstdFlags | log.LUTC
	defaultLogLevel  = LogWarning
)

var defaultLogger = NewDefaultLogger(defaultLogLevel)

func NewDefaultLogger(level LogLevel) *sampleLogger {
	return &sampleLogger{
		DEBUG: log.New(os.Stdout, fmt.Sprintf("DEBUG %s ", defaultLogPrefix), defaultLogFlag),
		ERROR: log.New(os.Stdout, fmt.Sprintf("ERROR %s ", defaultLogPrefix), defaultLogFlag),
		INFO:  log.New(os.Stdout, fmt.Sprintf("INFO %s ", defaultLogPrefix), defaultLogFlag),
		WARN:  log.New(os.Stdout, fmt.Sprintf("WARN %s ", defaultLogPrefix), defaultLogFlag),
		level: level,
	}
}

func (d *sampleLogger) Debug(ctx context.Context, module string, v ...interface{}) {
	if d.level <= LogDebug {
		d.DEBUG.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Debugf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= LogDebug {
		d.DEBUG.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Error(ctx context.Context, module string, v ...interface{}) {
	if d.level <= LogErr {
		d.ERROR.Output(2, fmt.Sprintln(v...))
	}

}

func (d *sampleLogger) Errorf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= LogErr {
		d.ERROR.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Info(ctx context.Context, module string, v ...interface{}) {
	if d.level >= LogInfo {
		d.INFO.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Infof(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level >= LogInfo {
		d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Warn(ctx context.Context, module string, v ...interface{}) {
	if d.level <= LogWarning {
		d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Warnf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= LogWarning {
		d.WARN.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Level() LogLevel {
	return d.level
}

func (d *sampleLogger) SetLevel(level LogLevel) {
	d.level = level
}
