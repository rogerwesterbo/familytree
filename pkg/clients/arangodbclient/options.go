package arangodbclient

import "time"

// Options holds configuration options for ArangoDB client
type Options struct {
	Endpoints  []string
	Username   string
	Password   string
	Database   string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

// OptionFunc is a function that configures Options
type OptionFunc func(*Options)

// WithEndpoints sets the ArangoDB endpoints
func WithEndpoints(endpoints []string) OptionFunc {
	return func(o *Options) {
		o.Endpoints = endpoints
	}
}

// WithUsername sets the ArangoDB username
func WithUsername(username string) OptionFunc {
	return func(o *Options) {
		o.Username = username
	}
}

// WithPassword sets the ArangoDB password
func WithPassword(password string) OptionFunc {
	return func(o *Options) {
		o.Password = password
	}
}

// WithDatabase sets the ArangoDB database name
func WithDatabase(database string) OptionFunc {
	return func(o *Options) {
		o.Database = database
	}
}

// WithTimeout sets the connection timeout
func WithTimeout(timeout time.Duration) OptionFunc {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(maxRetries int) OptionFunc {
	return func(o *Options) {
		o.MaxRetries = maxRetries
	}
}

// WithRetryDelay sets the delay between retries
func WithRetryDelay(retryDelay time.Duration) OptionFunc {
	return func(o *Options) {
		o.RetryDelay = retryDelay
	}
}

// DefaultOptions returns default options for ArangoDB client
func DefaultOptions() *Options {
	return &Options{
		Endpoints:  []string{"http://localhost:8529"},
		Username:   "root",
		Password:   "",
		Database:   "familytree",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// ApplyOptions applies option functions to Options
func ApplyOptions(opts *Options, optFuncs ...OptionFunc) *Options {
	for _, fn := range optFuncs {
		fn(opts)
	}
	return opts
}
