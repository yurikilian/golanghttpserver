package logger

var (
	Log = NewProvider().ProvideLog()
)
