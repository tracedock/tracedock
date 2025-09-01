package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var configContent = []byte(`
log:
  level: INFO

plugins:
  folders: [/etc/trackdock/plugins]

performance:
  memory_limiter:
    strategy: disk_dump
    max_consumption: 4096m

pipelines:
- name: main
  rules:
  - provider: eraser
    match:
      attributes:
        db.system: redis
        db.statement: ^HGETALL.*
      duration:
        lt: 100ms
`)

var configUnmarshaled = &Config{
	Log: ConfigLog{
		Level: "INFO",
	},
	Plugins: ConfigPlugins{
		Folders: []string{"/etc/trackdock/plugins"},
	},
	Performance: ConfigPerformance{
		MemoryLimiter: ConfigPerformanceMemoryLimiter{
			Strategy:       "disk_dump",
			MaxConsumption: "4096m",
		},
	},
	Pipelines: []ConfigPipeline{
		{
			Name: "main",
			Rules: []ConfigPipelineRules{
				{
					Provider: "eraser",
					Match: map[string]map[string]string{
						"attributes": {
							"db.system":    "redis",
							"db.statement": "^HGETALL.*",
						},
						"duration": {
							"lt": "100ms",
						},
					},
				},
			},
		},
	},
}

func Test_Config_Load(t *testing.T) {
	t.Cleanup(viper.Reset)

	t.Run("should return error when file doesn't exists", func(t *testing.T) {
		cfg := NewConfig()
		err := cfg.Load("non_existing_file.yaml")

		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("should load config when file exists", func(t *testing.T) {
		file, err := os.CreateTemp("", "tracedock_test_*.yaml")
		assert.NoError(t, err)

		file.Write(configContent)

		cfg := NewConfig()
		err = cfg.Load(file.Name())

		assert.NoError(t, err)
		assert.Equal(t, configUnmarshaled, cfg)
	})

	t.Run("should load default config when name isn't provided", func(t *testing.T) {
		cfg := NewConfig()
		err := cfg.Load("")

		assert.Error(t, err, os.ErrNotExist)
		assert.Equal(t, "/etc/tracedock/config.yaml", viper.ConfigFileUsed())
	})
}
