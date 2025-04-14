package log

import (
	"context"
	"fmt"
	"log"
	"os"
)

// @brief Level of the log message
//
// @since v0.7.0
type Level int

const (
	// DebugLevel is the lowest log level. It is used for debugging purposes.
	DebugLevel Level = iota
	// InfoLevel is used for informational messages.
	InfoLevel
	// WarningLevel is used for warning messages.
	WarningLevel
	// ErrLevel is used for error messages.
	ErrLevel
)

// @brief Logger interface,implement this interface to output logs to your logging component
//
// @since v0.7.0
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
	return DebugLevel
}

func (d *discardLogger) SetLevel(level Level) {}

type SampleLogger struct {
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

// @brief Creates a default logger with the specified log level
//
// @param level Log level. See Level for details.
//
// @since v0.7.0
func NewDefaultLogger(level Level) *SampleLogger {
	return &SampleLogger{
		DEBUG: log.New(os.Stdout, fmt.Sprintf("DEBUG %s ", defaultLogPrefix), defaultLogFlag),
		ERROR: log.New(os.Stdout, fmt.Sprintf("ERROR %s ", defaultLogPrefix), defaultLogFlag),
		INFO:  log.New(os.Stdout, fmt.Sprintf("INFO %s ", defaultLogPrefix), defaultLogFlag),
		WARN:  log.New(os.Stdout, fmt.Sprintf("WARN %s ", defaultLogPrefix), defaultLogFlag),
		level: level,
	}
}

func (d *SampleLogger) Debug(ctx context.Context, module string, v ...interface{}) {
	if d.level <= DebugLevel {
		_ = d.DEBUG.Output(2, fmt.Sprintln(v...))
	}
}

func (d *SampleLogger) Debugf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= DebugLevel {
		_ = d.DEBUG.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *SampleLogger) Error(ctx context.Context, module string, v ...interface{}) {
	if d.level <= ErrLevel {
		_ = d.ERROR.Output(2, fmt.Sprintln(v...))
	}
}

func (d *SampleLogger) Errorf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= ErrLevel {
		_ = d.ERROR.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *SampleLogger) Info(ctx context.Context, module string, v ...interface{}) {
	if d.level >= InfoLevel {
		_ = d.INFO.Output(2, fmt.Sprintln(v...))
	}
}

func (d *SampleLogger) Infof(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level >= InfoLevel {
		_ = d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *SampleLogger) Warn(ctx context.Context, module string, v ...interface{}) {
	if d.level <= WarningLevel {
		_ = d.WARN.Output(2, fmt.Sprintln(v...))
	}
}

func (d *SampleLogger) Warnf(ctx context.Context, module string, format string, v ...interface{}) {
	if d.level <= WarningLevel {
		_ = d.WARN.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *SampleLogger) Level() Level {
	return d.level
}

func (d *SampleLogger) SetLevel(level Level) {
	d.level = level
}
