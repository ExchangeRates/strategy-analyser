package config

type Config struct {
	Port          int    `toml:"port"`
	LogLevel      string `toml:"log_level"`
	RateResultUrl string `toml:"rate_result_url"`
}

func NewConfig() *Config {
	return &Config{
		Port:          8080,
		LogLevel:      "debug",
		RateResultUrl: "http://localhost:8081",
	}
}
