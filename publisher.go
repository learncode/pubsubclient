package pubsubclient

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

// Publisher contract to be returned to the consumer
type Publisher struct {
	client *pubSubClient
	topicsMap map[string]*pubsub.Topic
}

// PublisherConfig to be provided by the consumer.
type PublisherConfig struct {
	ProjectID string
}

// GetPublisher gives a publisher
func GetPublisher(config PublisherConfig) (*Publisher, error) {
	client, err := getClient(config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		client: client,
		topicsMap: make(map[string]*pubsub.Topic),
	}, nil
}

// Publish message to pubsub
func (publisher *Publisher) Publish(payload interface{}, topicNames []string) ([]string, error) {
	messages := make([]string, 0)
	data, err := json.Marshal(payload)
	if err != nil {
		return messages, err
	}
	message := &pubsub.Message{
		Data: data,
	}

	for _, topicName := range topicNames {
		topic, err := publisher.getTopic(topicName)
		if err != nil {
			return nil, err
		}
		response := topic.Publish(context.Background(), message)
		messageID, err := response.Get(context.Background())
		if err != nil {
			return messages, err
		}
		messages = append(messages, messageID)
	}

	return messages, nil
}

// getTopic gets the topic for the first time and then caches it locally
func (publisher *Publisher) getTopic(topicName string) (*pubsub.Topic, error) {
	if publisher.topicsMap[topicName] == nil {
		topic, err := publisher.client.createTopic(topicName)
		if err != nil {
			return nil, err
		}
		publisher.topicsMap[topicName] = topic
	}

	return publisher.topicsMap[topicName], nil
}
