package rapidpro

import (
	"time"

	"github.com/zewolfe/hermes/internal/log"
)

const (
	DefaultTimeout = time.Duration(30 * time.Second)
	DefaultTrigger = "dracarys"
)

type Options func(opts *options)

type options struct {
	token   string
	timeout time.Duration
	logger  *log.Logger

	//TODO: Maybe shouldn't be here
	trigger string
}

// TODO: This should probably not be optional
func WithToken(token string) Options {
	return func(opts *options) {
		opts.token = token
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

func WithLogger(logger *log.Logger) Options {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithTrigger(trigger string) Options {
	return func(opts *options) {
		opts.trigger = trigger
	}
}

func newWithOptions(opts ...Options) options {
	o := options{
		timeout: DefaultTimeout,
		trigger: DefaultTrigger,
	}

	for _, option := range opts {
		option(&o)
	}

	if o.logger == nil {
		o.logger = log.NewStdoutLogger()
	}

	return o
}
