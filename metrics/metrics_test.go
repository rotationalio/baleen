package metrics_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/metrics"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer(t *testing.T) {
	// Ensure that we can setup the monitoring server and execute metrics.
	conf := config.MonitoringConfig{
		Enabled:  true,
		BindAddr: "127.0.0.1:48489",
		NodeID:   "testing-42",
	}

	err := metrics.Serve(conf)
	require.NoError(t, err, "could not serve the metrics server")

	// Sleep a couple of milliseconds to ensure the metrics server is up
	time.Sleep(50 * time.Millisecond)

	// Collect some metrics
	metrics.Documents.WithLabelValues(conf.NodeID, "test").Inc()
	metrics.Subscriptions.Add(1)
	metrics.Subscriptions.Add(3)

	// Attempt to collect the metrics
	rep, err := http.Get("http://127.0.0.1:48489/metrics")
	require.NoError(t, err, "could not make http request to metrics server")
	require.Equal(t, http.StatusOK, rep.StatusCode)

	err = metrics.Shutdown(context.Background())
	require.NoError(t, err, "could not shutdown the metrics server")
}
