package logger

type Content interface {
	Msg(string)
	Msgf(string, ...any)
	Str(string, string) Content
	Err(error) Content
}

type Level interface {
	Info() Content
	Warn() Content
	Error() Content
	Fatal() Content
}

type Log interface {
	Level
}
