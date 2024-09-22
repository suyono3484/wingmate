package wingmate

import (
	"io"
	"time"

	"gitea.suyono.dev/suyono/wingmate/logger"
	"github.com/rs/zerolog"
)

const (
	timeTag = "time"
)

var (
	w *wrapper
)

type wrapper struct {
	log zerolog.Logger
}

func NewLog(wc io.WriteCloser) error {
	w = &wrapper{
		log: zerolog.New(wc),
	}
	return nil
}

func Log() logger.Log {
	if w == nil {
		panic("nil internal logger")
	}
	return w
}

func (w *wrapper) Info() logger.Content {
	return (*eventWrapper)(w.log.Info().Time(timeTag, time.Now()))
}

func (w *wrapper) Warn() logger.Content {
	return (*eventWrapper)(w.log.Warn().Time(timeTag, time.Now()))
}

func (w *wrapper) Error() logger.Content {
	return (*eventWrapper)(w.log.Error().Time(timeTag, time.Now()))
}

func (w *wrapper) Fatal() logger.Content {
	return (*eventWrapper)(w.log.Fatal().Time(timeTag, time.Now()))
}

type eventWrapper zerolog.Event

func (w *eventWrapper) Msg(msg string) {
	(*zerolog.Event)(w).Msg(msg)
}

func (w *eventWrapper) Msgf(format string, data ...any) {
	(*zerolog.Event)(w).Msgf(format, data...)
}

func (w *eventWrapper) Str(key, value string) logger.Content {
	rv := (*zerolog.Event)(w).Str(key, value)
	return (*eventWrapper)(rv)
}

func (w *eventWrapper) Err(err error) logger.Content {
	rv := (*zerolog.Event)(w).Err(err)
	return (*eventWrapper)(rv)
}
