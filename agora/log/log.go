package log

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrLevel
	UnknownLevel
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

	Level() Level
	SetLevel(level Level)
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

func (d *discardLogger) Level() Level {
	return UnknownLevel
}

func (d *discardLogger) SetLevel(level Level) {}

type sampleLogger struct {
	DEBUG *log.Logger
	ERROR *log.Logger
	INFO  *log.Logger
	WARN  *log.Logger

	level Level
}

const (
	defaultLogPrefix = "agora-rest-client-go"
	defaultLogFlag   = log.Lshortfile | log.LstdFlags | log.LUTC
	defaultLogLevel  = WarningLevel
)

var DefaultLogger = NewDefaultLogger(defaultLogLevel)

func NewDefaultLogger(level Level) *sampleLogger {
	return &sampleLogger{
		DEBUG: log.New(os.Stdout, fmt.Sprintf("DEBUG %s ", defaultLogPrefix), defaultLogFlag),
		ERROR: log.New(os.Stdout, fmt.Sprintf("ERROR %s ", defaultLogPrefix), defaultLogFlag),
		INFO:  log.New(os.Stdout, fmt.Sprintf("INFO %s ", defaultLogPrefix), defaultLogFlag),
		WARN:  log.New(os.Stdout, fmt.Sprintf("WARN %s ", defaultLogPrefix), defaultLogFlag),
		level: level,
	}
}

func (d *sampleLogger) Debug(ctx context.Context, module string, v ...interface{}) {
	if d.level <= DebugLevel {
		_ = d.DEBUG.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Debugf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= DebugLevel {
		_ = d.DEBUG.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Error(ctx context.Context, module string, v ...interface{}) {
	if d.level <= ErrLevel {
		_ = d.ERROR.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Errorf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= ErrLevel {
		_ = d.ERROR.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Info(ctx context.Context, module string, v ...interface{}) {
	if d.level >= InfoLevel {
		_ = d.INFO.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Infof(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level >= InfoLevel {
		_ = d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Warn(ctx context.Context, module string, v ...interface{}) {
	if d.level <= WarningLevel {
		_ = d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *sampleLogger) Warnf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= WarningLevel {
		_ = d.WARN.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *sampleLogger) Level() Level {
	return d.level
}

func (d *sampleLogger) SetLevel(level Level) {
	d.level = level
}
