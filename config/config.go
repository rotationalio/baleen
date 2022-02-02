package config

type Config struct {
	AWS     AWSConfig   `split_words:"true"`
	Kafka   KafkaConfig `split_words:"true"`
	Testing bool        `split_words:"true" default:"false"`
}

type AWSConfig struct {
	Enabled bool   `split_words:"true" default:"true"`
	Region  string `split_words:"true"`
	Bucket  string `split_words:"true"`
}

type KafkaConfig struct {
	Enabled bool   `split_words:"true" default:"false"`
	URL     string `split_words:"true"`
}
