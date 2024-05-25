package interpreter

type Config struct {
	ErrorReporting int64
}

func NewDevConfig() *Config {
	return &Config{ErrorReporting: E_ALL}
}

func NewProdConfig() *Config {
	return &Config{ErrorReporting: 0}
}
