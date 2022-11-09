package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rotationalio/baleen/config"
	"github.com/rs/zerolog/log"
)

const NamespaceBaleen = "baleen"

var (
	// Baleen specific collectors
	Subscriptions prometheus.Gauge
	FeedSyncs     *prometheus.CounterVec
	FeedItems     *prometheus.CounterVec
	Documents     *prometheus.CounterVec
)

// Internal package variables for serving the collectors to the Prometheus scraper.
var (
	srv   *http.Server
	cfg   config.MonitoringConfig
	setup sync.Once
	mu    sync.Mutex
	err   error
)

func Serve(conf config.MonitoringConfig) error {
	// Guard against concurrent Serve and Shutdown
	mu.Lock()
	defer mu.Unlock()

	// Ensure that the initialization of the metrics and the server occurs only once.
	setup.Do(func() {
		// Register the collectors
		cfg = conf
		if err = registerCollectors(); err != nil {
			return
		}

		// If not enabled, simply return here
		if !cfg.Enabled {
			return
		}

		// Setup the prometheus handler and collectors server.
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv = &http.Server{
			Addr:         cfg.BindAddr,
			Handler:      mux,
			ErrorLog:     nil,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		// Serve the metrics server in its own go routine
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error().Err(err).Msg("metrics server shutdown prematurely")
			}
		}()

		log.Info().Str("addr", fmt.Sprintf("http://%s/metrics", conf.BindAddr)).Msg("metrics server started and ready for prometheus collector")
	})

	return err
}

// Shutdown the prometheus metrics collectors server and reset the package. This method
// should be called at least once by outside callers before the process shuts down to
// ensure that system resources are cleaned up correctly.
func Shutdown(ctx context.Context) error {
	// Guard against concurrent Serve and Shutdown
	mu.Lock()
	defer mu.Unlock()

	// If we're already shutdown don't panic
	if srv == nil {
		return nil
	}

	// Ensure that no matter what happens we reset the package so it can be served again.
	defer func() {
		srv = nil
		cfg = config.MonitoringConfig{}
		err = nil
		setup = sync.Once{}
	}()

	// Ensure there is a shutdown deadline so we don't block forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Initializes and registers the metric collectors in Prometheus. This function should
// only be called once from the Serve function. All new metrics must be defined in this
// function so that they can be used.
func registerCollectors() (err error) {
	// Track all collectors to make it easier to register them at the end of this
	// function. When adding new collectors make sure to increase the capacity.
	collectors := make([]prometheus.Collector, 0, 8)

	// Baleen Collectors
	Subscriptions = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceBaleen,
		Name:      "subscriptions",
		Help:      "the number of subscriptions currently handled by the node",
	})
	collectors = append(collectors, Subscriptions)

	FeedSyncs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceBaleen,
		Name:      "feed_syncs",
		Help:      "the number of times a feed sync has occurred",
	}, []string{"node", "status_code"})
	collectors = append(collectors, FeedSyncs)

	FeedItems = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceBaleen,
		Name:      "feed_items",
		Help:      "the number of feed times discovered across all feed syncs",
	}, []string{"node", "feed_id"})
	collectors = append(collectors, FeedItems)

	Documents = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceBaleen,
		Name:      "documents",
		Help:      "the number of documents fetched",
	}, []string{"node", "status_code"})
	collectors = append(collectors, Documents)

	// Register all the collectors
	for _, collector := range collectors {
		if err = prometheus.Register(collector); err != nil {
			log.Debug().Err(err).Msg("could not register collector")
			return err
		}
	}

	return nil
}
