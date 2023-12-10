package wingmate

type CronTimeType int

const (
	Any CronTimeType = iota
	Exact
	MultipleOccurrence

	EnvPrefix = "WINGMATE"
)
