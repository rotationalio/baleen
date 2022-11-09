package config_test

import (
	"os"
	"testing"

	"github.com/rotationalio/baleen/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"BALEEN_LOG_LEVEL":                "debug",
	"BALEEN_CONSOLE_LOG":              "true",
	"BALEEN_AWS_ENABLED":              "false",
	"BALEEN_KAFKA_ENABLED":            "false",
	"BALEEN_MONITORING_ENABLED":       "true",
	"BALEEN_MONITORING_BIND_ADDR":     ":8889",
	"BALEEN_MONITORING_NODE_ID":       "test1234",
	"BALEEN_PUBLISHER_ENSIGN_ENABLED": "true",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after the test is complete.
	t.Cleanup(cleanupEnv())
	setEnv()

	conf, err := config.New()
	require.NoError(t, err, "could not process configuration from the environment")
	require.False(t, conf.IsZero(), "processed config should not be zero valued")

	// Ensure configuration is correctly set from the environment
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.True(t, conf.Monitoring.Enabled)
	require.Equal(t, testEnv["BALEEN_MONITORING_BIND_ADDR"], conf.Monitoring.BindAddr)
	require.Equal(t, testEnv["BALEEN_MONITORING_NODE_ID"], conf.Monitoring.NodeID)
}

// Returns the current environment for the specified keys, or if no keys are specified
// then it returns the current environment for all keys in the testEnv variable.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv variable. If no keys are specified,
// then this function sets all environment variables from the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}

// Cleanup helper function that can be run when the tests are complete to reset the
// environment back to its previous state before the test was run.
func cleanupEnv(keys ...string) func() {
	prevEnv := curEnv(keys...)
	return func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
