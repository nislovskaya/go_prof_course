package configs

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
