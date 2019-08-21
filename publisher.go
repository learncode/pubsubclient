package pubsubclient

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

// Publisher contract to be returned to the consumer
type Publisher struct {
	client    *pubSubClient
	topicsMap map[string]*pubsub.Topic
}

// PublisherConfig to be provided by the consumer.
type PublisherConfig struct {
	ProjectID string
}

// StopAll method will be used to publish all the messages and stop all
// the goroutines that handle publish. There by gauranteeing the messages
// are flushed out from the internal queue to google pub/sub
func (publisher *Publisher) StopAll() {
	for _, topic := range publisher.topicsMap {
		topic.Stop()
	}

	publisher.topicsMap = map[string]*pubsub.Topic{}
}

// GetPublisher gives a publisher
func GetPublisher(ctx context.Context, config PublisherConfig) (*Publisher, error) {
	client, err := getClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		client:    client,
		topicsMap: make(map[string]*pubsub.Topic),
	}, nil
}

// Publish message to pubsub
func (publisher *Publisher) Publish(ctx context.Context, payload interface{}, topicNames []string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	message := &pubsub.Message{
		Data: data,
	}

	for _, topicName := range topicNames {
		topic, err := publisher.getTopic(ctx, topicName)
		if err != nil {
			return err
		}
		response := topic.Publish(ctx, message)
		_, err = response.Get(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// getTopic gets the topic for the first time and then caches it locally
func (publisher *Publisher) getTopic(ctx context.Context, topicName string) (*pubsub.Topic, error) {
	if publisher.topicsMap[topicName] == nil {
		topic, err := publisher.client.createTopic(ctx, topicName)
		if err != nil {
			return nil, err
		}
		publisher.topicsMap[topicName] = topic
	}

	return publisher.topicsMap[topicName], nil
}
