package kit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.InterfaceMarshalFunc = Marshal
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFieldName = "ts"
	zerolog.TimestampFunc = Now
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000-0700"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = "panic"
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = "fatal"
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = "error"
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = "warning"
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = "info"
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = "debug"
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = "trace"

	FormatterText = "plain"
	FormatterJson = "json"
)

// ErrorHook allows specifying a hook for all logged errors
type ErrorHook interface {
	Error(err error)
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level   string // Level logging level
	Format  string // Format (plain, json)
	Context bool   // Context if true, context params are part of logging
	Service bool   // Service if true, service params are part of logging
	Caller  bool   // Caller if true, caller info is part of logging
}

type Logger struct {
	core zerolog.Logger
	cfg  *LogConfig
	hook ErrorHook
}

func InitLogger(cfg *LogConfig) *Logger {
	logger := &Logger{
		cfg: cfg,
	}
	logger.Init(cfg)
	return logger
}

func (l *Logger) SetErrorHook(h ErrorHook) {
	l.hook = h
}

func (l *Logger) Init(cfg *LogConfig) {
	l.cfg = cfg

	l.core = zerolog.New(os.Stdout)

	if cfg.Format == FormatterText {
		l.core = l.core.Output(zerolog.ConsoleWriter{
			Out:             os.Stderr,
			NoColor:         true,
			TimeFormat:      "2006-01-02T15:04:05.000-0700",
			TimeLocation:    nil,
			PartsExclude:    nil,
			FieldsOrder:     nil,
			FieldsExclude:   nil,
			FormatTimestamp: nil,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("|%-6s|", i))
			},
			FormatCaller: nil,
			FormatMessage: func(i interface{}) string {
				if i == nil {
					return ""
				}
				if s, ok := i.(string); ok && s == "" {
					return s
				}
				return fmt.Sprint(i)
			},
			FormatFieldName: func(i interface{}) string {
				if i == nil {
					return ""
				}
				if s, ok := i.(string); ok && s == "" {
					return s
				}
				return fmt.Sprintf("| %s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				if i == nil {
					return ""
				}
				if s, ok := i.(string); ok && s == "" {
					return s
				}
				return fmt.Sprint(i)
			},
			FormatErrFieldName: func(i interface{}) string {
				if i == nil {
					return ""
				}
				if s, ok := i.(string); ok && s == "" {
					return s
				}
				return fmt.Sprintf("| %s:", i)
			},
			FormatErrFieldValue:   nil,
			FormatPartValueByName: nil,
			FormatExtra:           nil,
			FormatPrepare:         nil,
		})

	} else {
	}

	lCtx := l.core.With().Timestamp()
	if cfg.Caller {
		lCtx = lCtx.Caller()
	}
	l.core = lCtx.Logger().Level(l.toZlLevel(cfg.Level))

}

func (l *Logger) GetLogger() *zerolog.Logger {
	return &l.core
}

func (l *Logger) SetLevel(level string) {
	l.cfg.Level = level
	l.Init(l.cfg)
}

func (l *Logger) toZlLevel(lv string) zerolog.Level {
	switch strings.ToLower(lv) {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	default:
		return zerolog.Disabled
	}
}

type CLoggerFunc func() CLogger

// CLogger provides structured logging abilities
// !!!! Not thread safe. Don't share one CLogger instance through multiple goroutines
type CLogger interface {
	// C adds request context to the log.
	// NOTE: Do not add context when logging errors. The context of where the error
	// occurred is usually different from where the log is invoked. Adding it in both
	// places may cause duplication and confusion.
	C(ctx context.Context) CLogger
	// F adds structured fields to the log.
	F(fields KV) CLogger
	// E attaches an error to the log entry.
	E(err error) CLogger
	// St appends a stack trace to the log entry, if an error is already set.
	St() CLogger
	// Cmp sets the component field (e.g., the subsystem or module name).
	Cmp(c string) CLogger
	// Mth sets the method or function name.
	Mth(m string) CLogger
	// Pr sets the protocol used (e.g., HTTP, gRPC).
	Pr(m string) CLogger
	// Srv sets the unique service code (useful in multi service logs).
	Srv(s string) CLogger
	// Node sets the node or instance identifier (e.g., hostname, container ID).
	Node(n string) CLogger
	// Inf logs at Info level with unformatted arguments.
	Inf(args ...interface{}) CLogger
	// InfF logs at Info level with a formatted message.
	InfF(format string, args ...interface{}) CLogger
	// Err logs at Error level with unformatted arguments.
	Err(args ...interface{}) CLogger
	// ErrF logs at Error level with a formatted message.
	ErrF(format string, args ...interface{}) CLogger
	// Dbg logs at Debug level with unformatted arguments.
	Dbg(args ...interface{}) CLogger
	// DbgF logs at Debug level with a formatted message.
	DbgF(format string, args ...interface{}) CLogger
	// Trc logs at Trace level with unformatted arguments.
	Trc(args ...interface{}) CLogger
	// TrcF logs at Trace level with a formatted message.
	TrcF(format string, args ...interface{}) CLogger
	// TrcObj logs objects only if the log level is set to Trace.
	// This avoids unnecessary marshaling overhead when Trace is disabled.
	// Note: Only exported fields will be serialized (per json.Marshal behavior).
	TrcObj(format string, args ...interface{}) CLogger
	// Warn logs at Warn level with unformatted arguments.
	Warn(args ...interface{}) CLogger
	// WarnF logs at Warn level with a formatted message.
	WarnF(format string, args ...interface{}) CLogger
	// Fatal logs at Fatal level and then exits.
	Fatal(args ...interface{}) CLogger
	// FatalF logs at Fatal level with a formatted message and then exits.
	FatalF(format string, args ...interface{}) CLogger
	// Clone creates a copy of the logger.
	// Always use Clone when passing CLogger between goroutines to avoid a shared state.
	Clone() CLogger
	// Printf logs a formatted debug message (alias for DbgF).
	Printf(format string, args ...interface{})
	// PrintfErr logs a formatted error message (alias for ErrF).
	PrintfErr(format string, args ...interface{})
	// Write implements io.Writer and logs the written bytes at Trace level.
	Write(p []byte) (n int, err error)
}

func L(logger *Logger) CLogger {
	return &clogger{
		logger: logger,
		core:   logger.core,
		fields: make(map[string]any),
		hook:   logger.hook,
	}
}

type clogger struct {
	logger *Logger
	core   zerolog.Logger
	err    error
	fields map[string]any
	hook   ErrorHook
	stack  string
}

// Clone must be used when passing CLogger between goroutines
func (cl *clogger) Clone() CLogger {

	clone := &clogger{
		core:   cl.core,
		err:    cl.err,
		fields: make(map[string]any),
		hook:   cl.hook,
	}

	for k, v := range cl.fields {
		clone.fields[k] = v
	}

	return clone
}

func (cl *clogger) C(ctx context.Context) CLogger {
	if !cl.logger.cfg.Context {
		return cl
	}
	if r, ok := Request(ctx); ok && r != nil {
		if rid := r.GetRequestId(); rid != "" {
			cl.fields["ctx.rid"] = rid
		}
		if un := r.GetUsername(); un != "" {
			cl.fields["ctx.un"] = un
		}
		if sid := r.GetSessionId(); sid != "" {
			cl.fields["ctx.sid"] = sid
		}
	}
	return cl
}

func (cl *clogger) F(fields KV) CLogger {
	for k, v := range fields {
		cl.fields[k] = v
	}
	return cl
}

func (cl *clogger) E(err error) CLogger {
	// if err is AppErr, log error code as a separate field
	if appErr, ok := IsAppErr(err); ok {
		cl.fields["err.code"] = appErr.Code()
		cl.fields["err.msg"] = appErr.Message()
		cl.fields["err.type"] = appErr.Type()

		for k, v := range appErr.Fields() {
			cl.fields[k] = v
		}
	}
	cl.err = err
	return cl
}

func (cl *clogger) St() CLogger {
	if cl.err != nil {
		// if err is AppErr take stack from an error itself, otherwise build stack right here
		if appErr, ok := IsAppErr(cl.err); ok {
			cl.stack = appErr.WithStack()
		} else {
			buf := make([]byte, 1<<16)
			s := runtime.Stack(buf, false)
			cl.stack = string(buf[0:s])
		}
	}
	return cl
}

func (cl *clogger) Srv(s string) CLogger {
	if !cl.logger.cfg.Service {
		return cl
	}
	cl.fields["call.svc"] = s
	return cl
}

func (cl *clogger) Node(n string) CLogger {
	if !cl.logger.cfg.Service {
		return cl
	}
	cl.fields["call.node"] = n
	return cl
}

func (cl *clogger) Cmp(c string) CLogger {
	cl.fields["call.cmp"] = c
	return cl
}

func (cl *clogger) Pr(c string) CLogger {
	cl.fields["call.pr"] = c
	return cl
}

func (cl *clogger) Mth(m string) CLogger {
	cl.fields["call.mth"] = m
	return cl
}

func (cl *clogger) Err(args ...interface{}) CLogger {

	l := cl.core.With().Fields(cl.fields).Err(cl.err).Logger()

	var msgSb strings.Builder

	if cl.stack != "" {
		msgSb.WriteString(cl.stack)
	}

	msgSb.WriteString(fmt.Sprint(args...))

	l.Error().Msg(msgSb.String())

	cl.fireHook()

	return cl
}

func (cl *clogger) ErrF(format string, args ...interface{}) CLogger {

	l := cl.core.With().Fields(cl.fields).Err(cl.err).Logger()
	l.Error().Msg(fmt.Sprintf(format, args...))

	cl.fireHook()

	return cl
}

func (cl *clogger) Inf(args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Info().Msg(fmt.Sprint(args...))
	return cl
}

func (cl *clogger) InfF(format string, args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Info().Msg(fmt.Sprintf(format, args...))
	return cl
}

func (cl *clogger) Warn(args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Warn().Msg(fmt.Sprint(args...))
	return cl
}

func (cl *clogger) WarnF(format string, args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Warn().Msg(fmt.Sprintf(format, args...))
	return cl
}

func (cl *clogger) Dbg(args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Debug().Msg(fmt.Sprint(args...))
	return cl
}

func (cl *clogger) DbgF(format string, args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Debug().Msg(fmt.Sprintf(format, args...))
	return cl
}

func (cl *clogger) Trc(args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Trace().Msg(fmt.Sprint(args...))
	return cl
}

func (cl *clogger) TrcF(format string, args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Trace().Msg(fmt.Sprintf(format, args...))
	return cl
}

func (cl *clogger) TrcObj(format string, args ...interface{}) CLogger {
	if cl.logger.cfg.Level == TraceLevel {
		var argsJs []interface{}
		for _, a := range args {
			if a != nil {
				js, _ := Marshal(a)
				argsJs = append(argsJs, string(js))
			}
		}
		return cl.TrcF(format, argsJs...)
	}
	return cl
}

func (cl *clogger) Fatal(args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Fatal().Msg(fmt.Sprint(args...))

	cl.fireHook()

	return cl
}

func (cl *clogger) FatalF(format string, args ...interface{}) CLogger {
	l := cl.core.With().Fields(cl.fields).Logger()
	l.Fatal().Msg(fmt.Sprintf(format, args...))

	cl.fireHook()

	return cl
}

func (cl *clogger) Printf(f string, args ...interface{}) {
	cl.DbgF(f, args...)
}

func (cl *clogger) PrintfErr(f string, args ...interface{}) {
	cl.ErrF(f, args...)
}

func (cl *clogger) Write(p []byte) (n int, err error) {
	cl.Trc(string(p))
	return len(p), nil
}

func (cl *clogger) fireHook() {
	if cl.hook != nil && cl.err != nil {
		cl.hook.Error(cl.err)
	}
}
