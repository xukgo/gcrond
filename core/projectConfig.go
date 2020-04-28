package core

type ProjectConfig struct {
	TimerTasks     []*TimerTaskConfig `yaml:"TimerTask"`
	RuleExecConfig []*RuleExecConfig  `yaml:"RuleExec"`
}
