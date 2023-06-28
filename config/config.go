package config

import (
	"errors"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/baleen/logger"
	"github.com/rotationalio/go-ensign"
	"github.com/rs/zerolog"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be BALEEN_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const prefix = "baleen"

// Config contains all of the configuration parameters for an Baleen service and is
// loaded from the environment or a configuration file with reasonable defaults for
// values that are omitted. The Config should be validated in preparation for running
// Baleen to ensure that all eventing operations work as expected.
// TODO: collect the config from a file instead of the environment.
type Config struct {
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog   bool                `split_words:"true" default:"false"`
	CloseTimeout time.Duration       `split_words:"true" default:"30s"`
	FeedSync     FeedSyncConfig      `split_words:"true"`
	PostFetch    PostFetchConfig     `split_words:"true"`
	Monitoring   MonitoringConfig
	Publisher    PublisherConfig
	Subscriber   SubscriberConfig
	processed    bool
}

type FeedSyncConfig struct {
	Enabled  bool          `default:"false"`
	Interval time.Duration `default:"1h"`
}

type PostFetchConfig struct {
	Enabled bool `default:"false"`
}

// MonitoringConfig maintains the parameters for the metrics server that the Prometheus
// scraper will fetch the configured observability metrics from.
type MonitoringConfig struct {
	Enabled  bool   `default:"true"`
	BindAddr string `split_words:"true" default:":1205"`
	NodeID   string `split_words:"true" required:"false"`
}

// Publisher Config defines the type of configuration to connect to the publisher with.
type PublisherConfig struct {
	Ensign EnsignConfig
	Kafka  KafkaConfig
}

// Subscriber Config defines the type of configuration to connect to the publisher with.
type SubscriberConfig struct {
	Ensign EnsignConfig
	Kafka  KafkaConfig
}

type EnsignConfig struct {
	Enabled      bool   `default:"true"`
	Endpoint     string `default:"ensign.rotational.app:443"`
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Insecure     bool   `default:"false"`
}

type KafkaConfig struct {
	Enabled        bool   `default:"false"`
	URL            string `split_words:"true"`
	Balancer       string `default:"LeastBytes"`
	TopicDocuments string `default:"documents"`
	TopicFeeds     string `default:"feeds"`
}

type AWSConfig struct {
	Enabled bool   `default:"false"`
	Region  string `split_words:"true"`
	Bucket  string `split_words:"true"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process(prefix, &conf); err != nil {
		return Config{}, err
	}

	// Post-process ensign config
	if conf.Publisher.Ensign.Enabled {
		conf.Publisher.Ensign.PostProcess()
	}

	if conf.Subscriber.Ensign.Enabled {
		conf.Subscriber.Ensign.PostProcess()
	}

	// Validate config-specific constraints
	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Parse and return the zerolog log level for configuring global logging.
func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

// A Config is zero-valued if it hasn't been processed by a file or the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Mark a manually constructed config as processed as long as its valid.
func (c Config) Mark() (Config, error) {
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

// Validate the entire config.
func (c Config) Validate() (err error) {
	if err = c.Publisher.Validate(); err != nil {
		return err
	}

	if err = c.Subscriber.Validate(); err != nil {
		return err
	}

	return nil
}

// Returns the Watermill RouterConfig=
func (c Config) RouterConfig() message.RouterConfig {
	return message.RouterConfig{
		CloseTimeout: c.CloseTimeout,
	}
}

func (c PublisherConfig) Validate() error {
	if !c.Ensign.Enabled && !c.Kafka.Enabled {
		return errors.New("invalid configuration: at least one publisher must be enabled")
	}

	if c.Kafka.Enabled {
		return c.Kafka.Validate()
	}

	return nil
}

func (c SubscriberConfig) Validate() error {
	if !c.Ensign.Enabled && !c.Kafka.Enabled {
		return errors.New("invalid configuration: at least one subscriber must be enabled")
	}

	if c.Kafka.Enabled {
		return c.Kafka.Validate()
	}

	return nil
}

func (c *EnsignConfig) PostProcess() {
	// Update ensign credentials with actual ensign environment variables.
	if c.ClientID == "" {
		c.ClientID = os.Getenv(ensign.EnvClientID)
	}

	if c.ClientSecret == "" {
		c.ClientSecret = os.Getenv(ensign.EnvClientSecret)
	}
}

// Validate the Kafka config.
func (c KafkaConfig) Validate() (err error) {
	if c.Enabled {
		if c.URL == "" {
			return errors.New("invalid configuration: kafka url must be specified")
		}

		if c.Balancer == "" {
			return errors.New("invalid configuration: kafka balancer must be specified")
		}

		if c.TopicDocuments == "" {
			return errors.New("invalid configuration: kafka topic for documents must be specified")
		}

		if c.TopicFeeds == "" {
			return errors.New("invalid configuration: kafka topic for feeds must be specified")
		}
	}

	return nil
}

// Validate the AWS config.
func (c AWSConfig) Validate() (err error) {
	if c.Enabled {
		if c.Region == "" {
			return errors.New("invalid configuration: AWS region must be specified")
		}

		if c.Bucket == "" {
			return errors.New("invalid configuration: AWS bucket must be specified")
		}
	}

	return nil
}
