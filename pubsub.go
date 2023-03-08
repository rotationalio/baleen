package baleen

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rotationalio/baleen/config"
	esdk "github.com/rotationalio/go-ensign"
	"github.com/rotationalio/watermill-ensign/pkg/ensign"
)

// Names of available topics
// TODO: how do we name our topics better and ensure there is a valid namespace?
const (
	TopicSubscriptions = "io.rotational.baleen/subscriptions"
	TopicFeeds         = "io.rotational.baleen/feeds"
	TopicDocuments     = "io.rotational.baleen/documents"
)

func CreatePublisher(conf config.PublisherConfig, logger watermill.LoggerAdapter) (message.Publisher, error) {
	if conf.Ensign.Enabled {
		return CreateEnsignPublisher(conf.Ensign, logger)
	}

	if conf.Kafka.Enabled {
		return CreateKafkaPublisher(conf.Kafka, logger)
	}

	return nil, errors.New("invalid configuration: no publisher enabled")
}

func CreateEnsignPublisher(conf config.EnsignConfig, logger watermill.LoggerAdapter) (message.Publisher, error) {
	// TODO: move the ensign config to the watermill-ensign library to avoid multi-import
	opts := ensign.PublisherConfig{
		EnsignConfig: &esdk.Options{
			Endpoint:     conf.Endpoint,
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			Insecure:     conf.Insecure,
		},
	}
	return ensign.NewPublisher(opts, logger)
}

func CreateKafkaPublisher(conf config.KafkaConfig, logger watermill.LoggerAdapter) (message.Publisher, error) {
	return nil, errors.New("not implemented yet")
}

func CreateSubscriber(conf config.SubscriberConfig, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	if conf.Ensign.Enabled {
		return CreateEnsignSubscriber(conf.Ensign, logger)
	}

	if conf.Kafka.Enabled {
		return CreateKafkaSubscriber(conf.Kafka, logger)
	}

	return nil, errors.New("invalid configuration: no subscriber enabled")
}

func CreateEnsignSubscriber(conf config.EnsignConfig, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	// TODO: move the ensign config to the watermill-ensign library to avoid multi-import
	opts := ensign.SubscriberConfig{
		EnsignConfig: &esdk.Options{
			Endpoint:     conf.Endpoint,
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			Insecure:     conf.Insecure,
		},
	}
	return ensign.NewSubscriber(opts, logger)
}

func CreateKafkaSubscriber(conf config.KafkaConfig, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	return nil, errors.New("not implemented yet")
}
