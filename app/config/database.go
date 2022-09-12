package config

type Database struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Username           string `yaml:"user"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
}
