package pubsubclient

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

// SubscribeMessageHandler that handles the message
type SubscribeMessageHandler func(chan *pubsub.Message)

// ErrorHandler that logs the error received while reading a message
type ErrorHandler func(error)

// SubscriberConfig subscriber config
type SubscriberConfig struct {
	ProjectID        string
	TopicName        string
	SubscriptionName string
	ErrorHandler     ErrorHandler
	Handle           SubscribeMessageHandler
}

// Subscriber subscribe to a topic and pass each message to the
// handler function
type Subscriber struct {
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	errorHandler ErrorHandler
	handle       SubscribeMessageHandler
	cancel       func()
}

// CreateSubscription creates a subscription
func CreateSubscription(ctx context.Context, config SubscriberConfig) (*Subscriber, error) {
	client, err := getClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	topic, err := client.createTopic(ctx, config.TopicName)
	if err != nil {
		return nil, err
	}
	subscription, err := client.createSubscription(ctx, config.SubscriptionName, topic)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		topic:        topic,
		subscription: subscription,
		errorHandler: config.ErrorHandler,
		handle:       config.Handle,
	}, nil
}

// Process will start pulling from the pubsub. The process accepts a waitgroup as
// it will be easier for us to orchestrate a use case where one application needs
// more than one subscriber
func (subscriber *Subscriber) Process(ctx context.Context, wg *sync.WaitGroup) {
	log.Printf("Starting a Subscriber on topic %s", subscriber.topic.String())
	output := make(chan *pubsub.Message)
	go func(subscriber *Subscriber, output chan *pubsub.Message) {
		defer close(output)

		ctx, subscriber.cancel = context.WithCancel(ctx)
		err := subscriber.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			output <- msg
		})

		if err != nil {
			// The wait group is stopped or marked done when an error is encountered
			subscriber.errorHandler(err)
			subscriber.stop()
			wg.Done()
		}
	}(subscriber, output)

	subscriber.handle(output)
}

// Stop the subscriber, closing the channel that was returned by Start.
func (subscriber *Subscriber) stop() {
	if subscriber.cancel != nil {
		log.Print("Stopped the subscriber")
		subscriber.cancel()
	}
}
