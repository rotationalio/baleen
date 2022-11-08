/*
Package baleen is the top level library of the baleen language ingestion service. This
library provides the primary components for running the service as a long running
background daemon including the main service itself, configuration and other utilities.
*/
package baleen

import (
	"context"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/logger"
	"github.com/rotationalio/watermill-ensign/pkg/ensign"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Initializes zerolog with our default logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Millisecond

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

// Baleen is essentially a wrapper for a watermill router that configures different
// event handlers depending on the context of the process. Calling Run() will start
// the Baleen service, which will handle incoming events and dispatch new events.
type Baleen struct {
	router     *message.Router
	conf       config.Config
	publisher  message.Publisher
	subscriber message.Subscriber
}

func New(conf config.Config) (svc *Baleen, err error) {
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Configure logging (will modify logging globally for all packages!)
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	svc = &Baleen{
		conf: conf,
	}

	var logger watermill.LoggerAdapter = logger.New()
	if svc.router, err = message.NewRouter(conf.RouterConfig(), logger); err != nil {
		return nil, err
	}

	// SignalsHandler will gracefully shutdown Router when SIGTERM is received.
	// You can also close the router by just calling `r.Close()`.
	svc.router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	svc.router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)

	// TODO: generalize the publisher and subscriber to anything
	// TODO: collect ensign configuration from the environment
	if svc.publisher, err = ensign.NewPublisher(ensign.PublisherConfig{}, logger); err != nil {
		return nil, err
	}

	if svc.subscriber, err = ensign.NewSubscriber(ensign.SubscriberConfig{}, logger); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Baleen) Run(ctx context.Context) error {
	return s.router.Run(ctx)
}

func (s *Baleen) Close() error {
	return s.router.Close()
}
