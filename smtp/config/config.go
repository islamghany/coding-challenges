package config

type Config struct {
	Host     string
	Port     string
	Hostname string // For SMTP greeting
}

func DefaultConfig() *Config {
	return &Config{
		Host:     "0.0.0.0",
		Port:     "2525",
		Hostname: "localhost",
	}
}

func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}
