package config

type ErrConfig struct {
	Errors   []error
	Warnings []error
}

func (e ErrConfig) Error() string {
	return "error occured while loading config"
}
