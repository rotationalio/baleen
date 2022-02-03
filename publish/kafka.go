package publish

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/store"
	"github.com/segmentio/kafka-go"
)

// A KafkaPublisher publishes objects to a Kafka topic.
type KafkaPublisher struct {
	conf     config.KafkaConfig
	writer   *kafka.Writer
	messages []kafka.Message
}

func New(conf config.KafkaConfig) (publisher *KafkaPublisher, err error) {
	writer := &kafka.Writer{
		Addr: kafka.TCP(conf.URL),
	}
	switch conf.Balancer {
	case "RoundRobin":
		writer.Balancer = &kafka.RoundRobin{}
	case "LeastBytes":
		writer.Balancer = &kafka.LeastBytes{}
	case "Hash":
		writer.Balancer = &kafka.Hash{}
	case "Murmur2":
		writer.Balancer = &kafka.Murmur2Balancer{}
	case "CRC32":
		writer.Balancer = &kafka.CRC32Balancer{}
	default:
		return nil, fmt.Errorf("unknown balancer specified: %s", conf.Balancer)
	}

	return &KafkaPublisher{
		conf:     conf,
		writer:   writer,
		messages: make([]kafka.Message, 0),
	}, nil
}

type KafkaDocument struct {
	FeedID       string `json:"feed_id"`
	LanguageCode string `json:"language_code"`
	Year         int    `json:"year"`
	Month        string `json:"month"`
	Day          int    `json:"day"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Content      string `json:"content"`
	Encoding     string `json:",omitempty"`
	Link         string `json:"link"`
}

// Write a Document to a Kafka message.
func (p *KafkaPublisher) WriteDocument(doc *store.Document) (err error) {
	// Need to encode the content as a base64 string for JSON serialization.
	kafkdaDoc := &KafkaDocument{
		FeedID:       doc.FeedID,
		LanguageCode: doc.LanguageCode,
		Year:         doc.Year,
		Month:        doc.Month,
		Day:          doc.Day,
		Title:        doc.Title,
		Description:  doc.Description,
		Content:      base64.RawStdEncoding.EncodeToString(doc.Content),
	}

	// Marshal into JSON
	var data []byte
	if data, err = json.Marshal(kafkdaDoc); err != nil {
		return err
	}

	p.messages = append(p.messages, kafka.Message{
		Topic: p.conf.TopicDocuments,
		Value: data,
	})
	return nil
}

// Write a Feed to a Kafka message.
func (p *KafkaPublisher) WriteFeed(feed *store.Feed) (err error) {
	// Marshal into JSON
	var data []byte
	if data, err = json.Marshal(feed); err != nil {
		return err
	}

	p.messages = append(p.messages, kafka.Message{
		Topic: p.conf.TopicFeeds,
		Value: data,
	})
	return nil
}

// Publish all the existing messages to Kafka.
func (p *KafkaPublisher) PublishMessages() (err error) {
	return p.writer.WriteMessages(context.Background(), p.messages...)
}
