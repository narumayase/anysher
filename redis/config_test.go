package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	os.Setenv("CACHE_ADDRESS", "localhost:6380")
	os.Setenv("CACHE_PASSWORD", "a_password")
	os.Setenv("CACHE_DATABASE", "1")

	expectedConfig := struct {
		name        string
		address     string
		password    string
		database    int
		expectedCfg Config
	}{
		name:     "Valid configuration",
		address:  "localhost:6380",
		password: "a_password",
		database: 1,
		expectedCfg: Config{
			cacheDatabase: 1,
			cachePassword: "a_password",
			cacheAddress:  "localhost:6380",
		},
	}
	cfg := NewConfiguration()
	assert.Equal(t, expectedConfig.expectedCfg, cfg)
}
