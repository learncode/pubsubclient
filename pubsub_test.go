package pubsubclient_test

import (
	pubsubclient "backend-components/gopubsubclient"
	"context"
	"log"
	"sync"
	"testing"

	"cloud.google.com/go/pubsub"
)

// If you want to run this locally before running this make sure gcloud emulator is up and running
const projectName = `<please-enter-your-project-name>`
const topicName = `<please-enter-a-topic-name>`
const subscriptionName = `<please-enter-a-subscription-name>`

func TestTopics(t *testing.T) {
	t.Run("Create a topic, subscribe to it, publishing a message to the topic the message gets proessed", func(t *testing.T) {
		publisherConfig := pubsubclient.PublisherConfig{
			ProjectID: projectName,
		}
		log.Printf("Publisher config: %+v", publisherConfig)
		ctx := context.Background()

		publisher, err := pubsubclient.GetPublisher(ctx, publisherConfig)
		log.Printf("Publisher: %+v", publisher)
		if err != nil {
			t.Fatalf("Error occured while creating a publisher, Err: %v", err)
		}

		subscriberConfig := pubsubclient.SubscriberConfig{
			ProjectID:        projectName,
			TopicName:        topicName,
			SubscriptionName: subscriptionName,
			ErrorHandler: func(err error) {
				log.Printf("Subscriber error: %v", err)
			},
			Handle: func(output chan *pubsub.Message) {
				for {
					pMsg := <-output
					log.Printf("Message: %+v", string(pMsg.Data))
					pMsg.Ack()
				}
			},
		}
		log.Printf("Subscriber config: %+v", subscriberConfig)
		subscriber, err := pubsubclient.CreateSubscription(ctx, subscriberConfig)
		if err != nil {
			t.Fatalf("Error occured while creating a subscriber, Err: %v", err)
		}

		log.Printf("Subscriber: %+v", subscriber)

		var wg sync.WaitGroup
		wg.Add(1)
		go subscriber.Process(ctx, &wg)
		publishMessages(ctx, publisher)
		wg.Wait()
	})
}

func publishMessages(ctx context.Context, publisher *pubsubclient.Publisher) {
	defer publisher.StopAll()
	for i := 0; i < 100; i++ {
		err := publisher.Publish(ctx, i, []string{topicName})
		if err != nil {
			log.Printf("Error occured while publishing the message, Err: %v", err)
		}
	}
}
