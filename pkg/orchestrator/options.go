package orchestrator

import (
	"time"

	"github.com/zewolfe/hermes/internal/log"
)

const (
	DefaultInterval = time.Minute * 5
	DefaultItemTTL  = time.Minute * 5
)

type Options func(opts *options)

type options struct {
	log              *log.Logger
	evictionInterval time.Duration
	ttl              time.Duration
}

func WithLogger(log *log.Logger) Options {
	return func(opts *options) {
		opts.log = log
	}
}

func WithEvictionInterval(interval time.Duration) Options {
	return func(opts *options) {
		opts.evictionInterval = interval
	}
}

func WithItemTTL(ttl time.Duration) Options {
	return func(opts *options) {
		opts.ttl = ttl
	}
}

func applyOptions(opts ...Options) options {
	o := options{
		log:              log.NewStdoutLogger(),
		evictionInterval: DefaultInterval,
		ttl:              DefaultItemTTL,
	}

	for _, option := range opts {
		option(&o)
	}

	return o
}
