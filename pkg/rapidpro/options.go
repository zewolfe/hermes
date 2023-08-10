package rapidpro

import (
	"time"

	"github.com/zewolfe/hermes/internal/log"
)

const (
	DefaultTimeout = time.Duration(30 * time.Second)
	DefaultTrigger = "dracarys"
)

type Options func(opts *options) error

type options struct {
	token   string
	timeout time.Duration
	logger  *log.Logger

	//TODO: Maybe shouldn't be here
	trigger string
}

// TODO: This should probably not be optional
func WithToken(token string) Options {
	return func(opts *options) error {
		opts.token = token

		return nil
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(opts *options) error {
		opts.timeout = timeout

		return nil
	}
}

func WithLogger(logger *log.Logger) Options {
	return func(opts *options) error {
		opts.logger = logger

		return nil
	}
}

func WithTrigger(trigger string) Options {
	return func(opts *options) error {
		opts.trigger = trigger

		return nil
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
