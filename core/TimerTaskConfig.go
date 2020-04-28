package core

type TimerTaskConfig struct {
	IntervalStr string   `yaml:"interval"`
	DoShellPath []string `yaml:"shellPath"`
	IntervalSec int64    `yaml:"-"`
}
