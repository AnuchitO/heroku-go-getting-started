package logger

import (
	"testing"
)

type mock struct {
	Env string
}

func (m mock) Getenv(name string) string {
	return m.Env
}

func TestLoadConfigLogDevelopmentEnv(t *testing.T) {
	mockEnv := mock{Env: "LOCAL"}
	config := loadConfigLog(mockEnv)

	if config.EncoderConfig.TimeKey != "T" {
		t.Errorf("Expected time key to be 'time', but got '%s'", config.EncoderConfig.TimeKey)
	}

	if config.DisableStacktrace {
		t.Error("Expected DisableStacktrace to be false, but it's true")
	}

	if len(config.OutputPaths) != 1 || config.OutputPaths[0] != "stderr" {
		t.Errorf("Expected OutputPaths to be ['log/ariskill.log'], but got %v", config.OutputPaths)
	}
}

func TestLoadConfigLogProductionEnv(t *testing.T) {
	mockEnv := mock{Env: "prod"}
	config := loadConfigLog(mockEnv)

	if config.EncoderConfig.TimeKey != "timestamp" {
		t.Errorf("Expected time key to be 'timestamp', but got '%s'", config.EncoderConfig.TimeKey)
	}

	if !config.DisableStacktrace {
		t.Errorf("Expected DisableStacktrace to be true, but got false")
	}

	// if len(config.OutputPaths) != 1 || config.OutputPaths[0] != "log/ariskill.log" {
	// 	t.Errorf("Expected OutputPaths to be ['log/ariskill.log'], but got %v", config.OutputPaths)
	// }
}
