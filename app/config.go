package app

import (
	"flag"
	"time"
)

// Config holds any kind of configuration that comes from the outside world and is necessary for running the application.
type Config struct {
	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	//
	// "development" is treated special: address types will use the loopback interface as default when none is defined.
	// This is useful when developing locally and listening on all interfaces requires elevated rights.
	Environment string `default:"production"`

	// Turns on some debug functionality: more verbose logs, exposed pprof, expvar and net trace in the debug server.
	Debug bool `split_words:"true"`

	// Defines the log format.
	// Valid values are: json, logfmt
	LogFormat string `split_words:"true" default:"json"`

	// Address of the debug server (configured by debug.addr flag)
	DebugAddr string `ignored:"true"`

	// Run the service in daemon mode instead of just running the job once.
	Daemon bool `ignored:"true"`

	// Schedule of the job when running in daemon mode.
	DaemonSchedule time.Duration `split_words:"true"`
}

// Flags configures a FlagSet.
//
// It still requires resolution (call to FlagSet.Parse) which is out of scope for this method.
func (c *Config) Flags(flags *flag.FlagSet) {
	flags.StringVar(&c.DebugAddr, "debug.addr", ":10000", "Debug and health check address")

	flags.BoolVar(&c.Daemon, "daemon", false, "Start as daemon instead of single run")
}
