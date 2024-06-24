package logger

import (
	"fmt"
	"github.com/onflow/flowkit/v2/output"
	"github.com/rs/zerolog"
	"io"
	"os"
)

type Logger struct {
	logger *zerolog.Logger
	cache  *CacheLogWriter
}

var _ output.Logger = &Logger{}

type Config struct {
	Verbose   bool
	LogFormat string // "text" or "json". Defaults to "json" if "logs" writer is used.
}

func NewLogger(config Config) *Logger {

	level := zerolog.InfoLevel
	if config.Verbose {
		level = zerolog.DebugLevel
	}
	zerolog.MessageFieldName = "msg"

	cacheWriter := NewCacheLogWriter()

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
		cacheWriter,
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &Logger{
		logger: &logger,
		cache:  cacheWriter,
	}
}

func (l *Logger) Zerolog() *zerolog.Logger {
	return l.logger
}

func (l *Logger) Debug(s string) {
	l.Debug(s)
}

func (l *Logger) Info(s string) {
	l.Info(s)
}

func (l *Logger) Error(s string) {
	l.Error(s)
}

func (l *Logger) StartProgress(s string) {
	l.Info(fmt.Sprintf("🏗️ %s", s))
}

func (l *Logger) StopProgress() {
	// noop
}

func (l *Logger) LogsHistory() []string {
	return l.cache.logs
}

type CacheLogWriter struct {
	logs []string
}

func NewCacheLogWriter() *CacheLogWriter {
	return &CacheLogWriter{
		logs: make([]string, 0),
	}
}

var _ io.Writer = &CacheLogWriter{}

func (c *CacheLogWriter) Write(p []byte) (n int, err error) {
	c.logs = append(c.logs, string(p))
	return len(p), nil
}

func NewTextWriter() zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("%-44s", i)
	}

	return writer
}
