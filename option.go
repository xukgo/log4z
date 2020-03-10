package log4z

const (
	DefaultTimeKey       = "T"
	DefaultLevelKey      = "level"
	DefaultNameKey       = "logger"
	DefaultCallerKey     = "line"
	DefaultMessageKey    = "msg"
	DefaultStacktraceKey = "stacktrace"
	DefaultTimeFormat    = "2006-01-02 15:04:05.000"
	DefaultCompress      = true
	DefaultCompressDelay = 0
)

type Option func(opts *Options)

type Options struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string
	TimeFormat    string
	Compress      bool
	CompressDelay int
}

// WithOptions accepts the whole options config.
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func WithTimeKey(key string) Option {
	return func(opts *Options) {
		opts.TimeKey = key
	}
}

func WithLevelKey(key string) Option {
	return func(opts *Options) {
		opts.LevelKey = key
	}
}
func WithNameKey(key string) Option {
	return func(opts *Options) {
		opts.NameKey = key
	}
}
func WithCallerKey(key string) Option {
	return func(opts *Options) {
		opts.CallerKey = key
	}
}
func WithMessageKey(key string) Option {
	return func(opts *Options) {
		opts.MessageKey = key
	}
}
func WithStacktraceKey(key string) Option {
	return func(opts *Options) {
		opts.StacktraceKey = key
	}
}
func WithTimeFormat(key string) Option {
	return func(opts *Options) {
		opts.TimeFormat = key
	}
}

func WithCompress(key bool) Option {
	return func(opts *Options) {
		opts.Compress = key
	}
}
func WithCompressDelay(key int) Option {
	return func(opts *Options) {
		opts.CompressDelay = key
	}
}
