package logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger *ZeroLogger
)

func init() {
	logger = &ZeroLogger{}
}

func New() *ZeroLogger {
	return logger
}

type ZeroLogger struct{}

var _ watermill.LoggerAdapter = &ZeroLogger{}

func (ZeroLogger) Error(msg string, err error, fields watermill.LogFields) {
	log.Error().Err(err).Fields(map[string]interface{}(fields)).Msg(msg)
}

func (ZeroLogger) Info(msg string, fields watermill.LogFields) {
	log.Info().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (ZeroLogger) Debug(msg string, fields watermill.LogFields) {
	log.Debug().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (ZeroLogger) Trace(msg string, fields watermill.LogFields) {
	log.Trace().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (ZeroLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &ZeroContext{log: log.With().Fields(map[string]interface{}(fields)).Logger()}
}

type ZeroContext struct {
	log zerolog.Logger
}

var _ watermill.LoggerAdapter = &ZeroContext{}

func (z ZeroContext) Error(msg string, err error, fields watermill.LogFields) {
	z.log.Error().Err(err).Fields(map[string]interface{}(fields)).Msg(msg)
}

func (z ZeroContext) Info(msg string, fields watermill.LogFields) {
	z.log.Info().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (z ZeroContext) Debug(msg string, fields watermill.LogFields) {
	z.log.Debug().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (z ZeroContext) Trace(msg string, fields watermill.LogFields) {
	z.log.Trace().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (z ZeroContext) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &ZeroContext{log: z.log.With().Fields(map[string]interface{}(fields)).Logger()}
}
