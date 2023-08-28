package wingmate

import (
	"gitea.suyono.dev/suyono/wingmate/logger"
	"github.com/rs/zerolog"
	"io"
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
	return w.log.Info()
}

func (w *wrapper) Warn() logger.Content {
	return w.log.Warn()
}

func (w *wrapper) Error() logger.Content {
	return w.log.Error()
}
