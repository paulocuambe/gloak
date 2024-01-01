package config

type ConfigErr struct {
	Errs []error
}

func (e ConfigErr) Error() string {
	return "config error"
}

type ConfigWarnigs struct {
	Warnings []error
}

func (e ConfigWarnigs) Error() string {
	return "config warnings"
}
