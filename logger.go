package wingmate

import (
	"io"

	"gitea.suyono.dev/suyono/wingmate/logger"
	"github.com/rs/zerolog"
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
	return (*eventWrapper)(w.log.Info())
}

func (w *wrapper) Warn() logger.Content {
	return (*eventWrapper)(w.log.Warn())
}

func (w *wrapper) Error() logger.Content {
	return (*eventWrapper)(w.log.Error())
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
