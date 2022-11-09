package config

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/baleen/logger"
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
	CloseTimeout time.Duration       `split_words:"true"`
	AWS          AWSConfig           `split_words:"true"`
	Kafka        KafkaConfig         `split_words:"true"`
	DBPath       string              `split_words:"true" default:"./db"`
	FixturesDir  string              `split_words:"true" default:"fixtures"`
	Testing      bool                `split_words:"true" default:"false"`
	processed    bool
}

type AWSConfig struct {
	Enabled bool   `split_words:"true" default:"true"`
	Region  string `split_words:"true"`
	Bucket  string `split_words:"true"`
}

type KafkaConfig struct {
	Enabled        bool   `split_words:"true" default:"false"`
	URL            string `split_words:"true"`
	Balancer       string `split_words:"true" default:"LeastBytes"`
	TopicDocuments string `split_words:"true" default:"documents"`
	TopicFeeds     string `split_words:"true" default:"feeds"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process(prefix, &conf); err != nil {
		return Config{}, err
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
	if err = c.AWS.Validate(); err != nil {
		return err
	}

	if err = c.Kafka.Validate(); err != nil {
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

// Validate the AWS config.
func (c AWSConfig) Validate() (err error) {
	if c.Enabled {
		if c.Region == "" {
			return fmt.Errorf("AWS region must be specified")
		}

		if c.Bucket == "" {
			return fmt.Errorf("AWS bucket must be specified")
		}
	}

	return nil
}

// Validate the Kafka config.
func (c KafkaConfig) Validate() (err error) {
	if c.Enabled {
		if c.URL == "" {
			return fmt.Errorf("kafka url must be specified")
		}

		if c.Balancer == "" {
			return fmt.Errorf("kafka balancer must be specified")
		}

		if c.TopicDocuments == "" {
			return fmt.Errorf("kafka topic for documents must be specified")
		}

		if c.TopicFeeds == "" {
			return fmt.Errorf("kafka topic for feeds must be specified")
		}
	}

	return nil
}
