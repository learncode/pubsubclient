package pubsubclient

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

type pubSubClient struct {
	psclient *pubsub.Client
}

func getClient(projectID string) (*pubSubClient, error) {
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Printf("Error when creating pubsub client. Err: %v", err)
		return nil, err
	}
	return &pubSubClient{psclient: client}, nil
}

// topicExists checks if a given topic exists
func (client *pubSubClient) topicExists(topicName string) (bool, error) {
	topic := client.psclient.Topic(topicName)
	return topic.Exists(context.Background())
}

// createTopic creates a topic if a topic name does not exist or returns one
// if it is already present
func (client *pubSubClient) createTopic(topicName string) (*pubsub.Topic, error) {
	topicExists, err := client.topicExists(topicName)
	if err != nil {
		log.Printf("Could not check if topic exists. Error: %+v", err)
		return nil, err
	}
	var topic *pubsub.Topic

	if !topicExists {
		topic, err = client.psclient.CreateTopic(context.Background(), topicName)
		if err != nil {
			log.Printf("Could not create topic. Err: %+v", err)
			return nil, err
		}
	} else {
		topic = client.psclient.Topic(topicName)
	}

	return topic, nil
}

// createSubscription creates the subscription to a topic
func (client *pubSubClient) createSubscription(subscriptionName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	subscription := client.psclient.Subscription(subscriptionName)

	subscriptionExists, err := subscription.Exists(context.Background())
	if err != nil {
		log.Printf("Could not check if subscription %s exists. Err: %v", subscriptionName, err)
		return nil, err
	}

	if !subscriptionExists {

		cfg := pubsub.SubscriptionConfig{
			Topic: topic,
			// The subscriber has a configurable, limited amount of time -- known as the ackDeadline -- to acknowledge
			// the outstanding message. Once the deadline passes, the message is no longer considered outstanding, and
			// Cloud Pub/Sub will attempt to redeliver the message.
			AckDeadline: 60 * time.Second,
		}

		subscription, err = client.psclient.CreateSubscription(context.Background(), subscriptionName, cfg)
		if err != nil {
			log.Printf("Could not create subscription %s. Err: %v", subscriptionName, err)
			return nil, err
		}
		subscription.ReceiveSettings = pubsub.ReceiveSettings{
			// This is the maximum amount of messages that are allowed to be processed by the callback function at a time.
			// Once this limit is reached, the client waits for messages to be acked or nacked by the callback before
			// requesting more messages from the server.
			MaxOutstandingMessages: 100,
			// This is the maximum amount of time that the client will extend a message's deadline. This value should be
			// set as high as messages are expected to be processed, plus some buffer.
			MaxExtension: 10 * time.Second,
		}
	}
	return subscription, nil
}

// subscriptionExists checks if a given subscription exists
func (client *pubSubClient) subscriptionExists(subscriptionName string) (bool, error) {
	subscription := client.psclient.Subscription(subscriptionName)
	return subscription.Exists(context.Background())
}

// deleteSubscription deletes a subscription
func (client *pubSubClient) deleteSubscription(subscriptionName string) error {
	return client.psclient.Subscription(subscriptionName).Delete(context.Background())
}

// listAllSubscription lists all subscriptions in the project
func (client *pubSubClient) listAllSubscription(topicName string) ([]string, error) {
	subscriptionNames := make([]string, 0)
	subscriptionIterator := client.psclient.Subscriptions(context.Background())
	for {
		item, err := subscriptionIterator.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("Could not list all topics. Error %v", err)
			return subscriptionNames, err
		}
		subscriptionNames = append(subscriptionNames, item.String())
	}
	return subscriptionNames, nil
}
