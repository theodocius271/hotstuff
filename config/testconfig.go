package config

import (
	"encoding/json"
	"os"
	"testing"
)

// WriteTestConfig 为测试目的临时写入配置文件
func WriteTestConfig(config *HotStuffConfig, t *testing.T) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "hotstuff-test-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	_, err = tmpFile.Write(configBytes)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// 设置环境变量，使配置从此文件读取
	os.Setenv("HOTSTUFF_CONFIG_PATH", tmpFile.Name())

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
		os.Unsetenv("HOTSTUFF_CONFIG_PATH")
	})
}
