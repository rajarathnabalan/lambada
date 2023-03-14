package lambada

// OutputMode represents the way the request's output will be handled.
// See the defined OutputMode cconstant to get details on available output modes and how they work.
type OutputMode int8

const (
	// Fully manual mode. Neither Content-Type or binary mode will be automatically activated on the responses.
	// Binary mode may be manually set on the response writer if needed.
	Manual OutputMode = -1

	// The Content-Type will automatically be set in the response, if not already provided.
	// The binary mode will not be activated automatically and must be set manually when required.
	//
	// This is the default mode for backward compatibility reasons.
	AutoContentType OutputMode = 0

	// Fully automatic mode: The Content-Type will be set in the response, if not already provided.
	// If the Content-Type indicated binary data, binary mode is also activated automatically.
	// Binary mode can still be forcefully manually enabled, but cannot be forced to disabled.
	Automatic OutputMode = 1
)

// A Logger interface. Provide a single Printf method, which is compatible with the standard library log package.
type Logger interface {
	Printf(fmt string, args ...interface{})
}

// NullLogger is a no-op implementation of Logger. All messages sent to it are ignored.
type NullLogger struct{}

func (l NullLogger) Printf(fmt string, args ...interface{}) {}

type options struct {
	requestLogger  Logger
	responseLogger Logger
	outputMode     OutputMode
	defaultBinary  bool
}

// newOptions creates a new options and applies opts.
// Prior applying opts, the new options are initialized with the zero value for all fields except the loggers, which
// are all initialized with a NullLogger.
func newOptions(opts ...Option) *options {
	o := &options{
		requestLogger:  NullLogger{},
		responseLogger: NullLogger{},
	}
	o.apply(opts...)
	return o
}

// apply applies opts to o, in order.
func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// An Option used to customize the handler behavior
type Option func(*options)

// WithOutputMode sets the handler's output mode.
// OutputMode may be one of Manual, AutoContentType or Automatic.
func WithOutputMode(outputMode OutputMode) Option {
	return func(o *options) {
		o.outputMode = outputMode
	}
}

// WithLogger sets all the handler's loggers to logger.
func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.requestLogger = logger
		o.responseLogger = logger
	}
}

// WithRequestLogger sets the handler's request logger.
func WithRequestLogger(logger Logger) Option {
	return func(o *options) {
		o.requestLogger = logger
	}
}

// WithResponseLogger sets the handler's response logger.
func WithResponseLogger(logger Logger) Option {
	return func(o *options) {
		o.responseLogger = logger
	}
}

// WithDefaultBinary enables or disable the default binary mode.
func WithDefaultBinary(defaultBinary bool) Option {
	return func(o *options) {
		o.defaultBinary = defaultBinary
	}
}
