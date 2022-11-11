package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AWS         AWSConfig   `split_words:"true"`
	Kafka       KafkaConfig `split_words:"true"`
	DBPath      string      `split_words:"true" default:"./db"`
	FixturesDir string      `split_words:"true" default:"fixtures"`
	Testing     bool        `split_words:"true" default:"false"`
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
	if err = envconfig.Process("baleen", &conf); err != nil {
		return Config{}, err
	}

	// Validate config-specific constraints
	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	return conf, nil
}

// Validate the entire config.
func (c Config) Validate() (err error) {
	if c.DBPath == "" {
		return fmt.Errorf("DBPath must be set")
	}

	if c.FixturesDir == "" {
		return fmt.Errorf("FixturesDir must be set")
	}

	if err = c.AWS.Validate(); err != nil {
		return err
	}

	if err = c.Kafka.Validate(); err != nil {
		return err
	}

	return nil
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
			return fmt.Errorf("Kafka URL must be specified")
		}

		if c.Balancer == "" {
			return fmt.Errorf("Kafka balancer must be specified")
		}

		if c.TopicDocuments == "" {
			return fmt.Errorf("Kafka topic for documents must be specified")
		}

		if c.TopicFeeds == "" {
			return fmt.Errorf("Kafka topic for feeds must be specified")
		}
	}

	return nil
}
