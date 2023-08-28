package logger

type Content interface {
	Msg(string)
	Msgf(string, ...any)
}

type Level interface {
	Info() Content
	Warn() Content
	Error() Content
}

type Log interface {
	Level
}
