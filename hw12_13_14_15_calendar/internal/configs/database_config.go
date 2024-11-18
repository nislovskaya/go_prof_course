package configs

import "fmt"

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (dc *DatabaseConfig) String() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		dc.Driver,
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
		dc.SSLMode,
	)
}
