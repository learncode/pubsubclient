package pubsubclient_test

import (
	pubsubclient "backend-components/gopubsubclient"
	"log"
	"sync"
	"testing"

	"cloud.google.com/go/pubsub"
)

// Before running this make sure gcloud emulator is up and running

const projectName = `<please-enter-your-project-name>`
const topicName = `<please-enter-a-topic-name>`
const subscriptionName = `<please-enter-a-subscription-name>`

func TestTopics(t *testing.T) {
	t.Run("Create a topic, subscribe to it, publishing a message to the topic the message gets proessed", func(t *testing.T) {
		publisherConfig := pubsubclient.PublisherConfig{
			ProjectID: projectName,
		}
		log.Printf("Publisher config: %+v", publisherConfig)

		publisher, err := pubsubclient.GetPublisher(publisherConfig)
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
		subscriber, err := pubsubclient.CreateSubscription(subscriberConfig)
		if err != nil {
			t.Fatalf("Error occured while creating a subscriber, Err: %v", err)
		}

		log.Printf("Subscriber: %+v", subscriber)

		var wg sync.WaitGroup
		wg.Add(1)
		go subscriber.Process(&wg)
		publishMessages(publisher)
		wg.Wait()
	})
}

func publishMessages(publisher *pubsubclient.Publisher) {
	for i := 0; i < 100; i++ {
		id, err := publisher.Publish(i, []string{topicName})
		if err != nil {
			log.Printf("Error occured while publishing the message, Err: %v", err)
		}
		log.Printf("MessageID: %v", id)
	}
}
