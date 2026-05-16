package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("NOSTALGIA_SCAN_PATH", "/scan")
	os.Setenv("NOSTALGIA_HOME_PATH", "/home")
	os.Setenv("NOSTALGIA_THUMB_TARGET_PATH", "/thumb")
	os.Setenv("NOSTALGIA_CONNECTION_STRING", "db_conn")
	defer os.Unsetenv("NOSTALGIA_SCAN_PATH")
	defer os.Unsetenv("NOSTALGIA_HOME_PATH")
	defer os.Unsetenv("NOSTALGIA_THUMB_TARGET_PATH")
	defer os.Unsetenv("NOSTALGIA_CONNECTION_STRING")

	cfg := LoadConfig()

	assert.Equal(t, "/scan", cfg.NOSTALGIA_SCAN_PATH)
	assert.Equal(t, "/home", cfg.NOSTALGIA_HOME_PATH)
	assert.Equal(t, "/thumb", cfg.NOSTALGIA_THUMB_TARGET_PATH)
	assert.Equal(t, "db_conn", cfg.NOSTALGIA_CONNECTION_STRING)
}
