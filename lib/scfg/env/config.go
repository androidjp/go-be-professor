package env

const DefaultPrefix = "KLION_"

type Config struct {
	Enable   bool     `json:"enable"`
	Prefixes []string `json:"prefixes"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:   true,
		Prefixes: []string{DefaultPrefix},
	}
}
