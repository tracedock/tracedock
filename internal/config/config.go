package config

import (
	"github.com/spf13/viper"
)

type ConfigLog struct {
	Level string
}

type ConfigPlugins struct {
	Folders []string
}

type ConfigPerformanceMemoryLimiter struct {
	Strategy       string `mapstructure:"strategy"`
	MaxConsumption string `mapstructure:"max_consumption"`
}

type ConfigPerformance struct {
	MemoryLimiter ConfigPerformanceMemoryLimiter `mapstructure:"memory_limiter"`
}

type ConfigPipelineRules struct {
	Provider string

	Match   map[string]map[string]string
	Missing map[string]map[string]string
	Set     map[string]string
}

type ConfigPipeline struct {
	Name  string
	Rules []ConfigPipelineRules
}

type Config struct {
	Log         ConfigLog
	Plugins     ConfigPlugins
	Performance ConfigPerformance
	Pipelines   []ConfigPipeline
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(file string) error {
	if file == "" {
		file = "/etc/tracedock/config.yaml"
	}

	c.setDefaults()

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) setDefaults() {
	viper.SetDefault("log.level", "INFO")
	viper.SetDefault("plugins.folders", []string{"/etc/trackdock/plugins"})
}
