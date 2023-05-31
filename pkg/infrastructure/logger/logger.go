package logger

import (
	"bufio"
	"bytes"
	"context"
	"os"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/tracer"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const traceID = "trace_id"

// Logger provides the logging functionality.
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger initializes Logger instance.
func NewLogger() *Logger {
	c := newLoggerConfig()
	logger, _ := c.Build()

	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// NewTestLogger initializes a test Logger instance that is useful for testing purpose.
func NewTestLogger() (*Logger, *bytes.Buffer, *bufio.Writer) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	c := newLoggerConfig()

	return &Logger{
		SugaredLogger: zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(c.EncoderConfig),
				zapcore.AddSync(writer),
				zapcore.DebugLevel,
			),
		).Sugar(),
	}, &buffer, writer
}

// DebugContext uses fmt.Sprint to construct and log a message with the `trace_id` found in the context.
func (logger *Logger) DebugContext(ctx context.Context, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Debug(args...)
}

// DebugfContext uses fmt.Sprintf to log a templated message with the `trace_id` found in the context.
func (logger *Logger) DebugfContext(ctx context.Context, template string, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Debugf(template, args...)
}

// ErrorContext uses fmt.Sprint to construct and log a message with the `trace_id` found in the context.
func (logger *Logger) ErrorContext(ctx context.Context, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Error(args...)
}

// ErrorfContext uses fmt.Sprintf to log a templated message with the `trace_id` found in the context.
func (logger *Logger) ErrorfContext(ctx context.Context, template string, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Errorf(template, args...)
}

// InfoContext uses fmt.Sprint to construct and log a message with the `trace_id` found in the context.
func (logger *Logger) InfoContext(ctx context.Context, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Info(args...)
}

// InfofContext uses fmt.Sprintf to log a templated message with the `trace_id` found in the context.
func (logger *Logger) InfofContext(ctx context.Context, template string, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Infof(template, args...)
}

// WarnContext uses fmt.Sprint to construct and log a message with the `trace_id` found in the context.
func (logger *Logger) WarnContext(ctx context.Context, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Warn(args...)
}

// WarnfContext uses fmt.Sprintf to log a templated message with the `trace_id` found in the context.
func (logger *Logger) WarnfContext(ctx context.Context, template string, args ...interface{}) {
	logger.With(traceID, tracer.GetTraceIDFromContext(ctx)).Warnf(template, args...)
}

func newLoggerConfig() zap.Config {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	c.EncoderConfig.TimeKey = ""
	c.EncoderConfig.CallerKey = ""

	if os.Getenv("APP_ENV") != "" && os.Getenv("APP_ENV") != "development" {
		c = zap.NewProductionConfig()
	}

	return c
}
