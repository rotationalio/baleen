package publish

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/store"
	"github.com/segmentio/kafka-go"
)

// A KafkaPublisher publishes Documents to a Kafka topic.
type KafkaPublisher struct {
	writer *kafka.Writer
}

func New(config config.KafkaConfig) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(config.URL),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) PublishDocuments(documents []*store.Document) error {
	var errs *multierror.Error

	// Convert the documents into Kafka messages
	messages := make([]kafka.Message, 0, len(documents))
	for _, doc := range documents {
		var data []byte
		var err error
		if data, err = json.Marshal(doc); err != nil {
			// TODO: Need some way to know which documents errored so we can publish them later.
			errs = multierror.Append(errs, err)
			continue
		}

		messages = append(messages, kafka.Message{
			Topic: "documents",
			Value: []byte(data),
		})
	}

	if err := p.writer.WriteMessages(context.Background(), messages...); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

func (p *KafkaPublisher) PublishFeeds(feeds []*store.Feed) error {
	var errs *multierror.Error

	// Convert the documents into Kafka messages
	messages := make([]kafka.Message, 0, len(feeds))
	for _, feed := range feeds {
		var data []byte
		var err error
		if data, err = json.Marshal(feed); err != nil {
			// TODO: Need some way to know which feeds errored so we can publish them later.
			errs = multierror.Append(errs, err)
			continue
		}

		messages = append(messages, kafka.Message{
			Topic: "feeds",
			Value: []byte(data),
		})
	}

	if err := p.writer.WriteMessages(context.Background(), messages...); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}
