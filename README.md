# Pubsub client

A client library for `Google PUB/SUB` written in `GO`

## Features

+ Create a topic.
+ Create a subscription.
+ Publish messages to a topic.
+ Use a pull subscriber to output individual topic messages.

## Creating the publisher

```go
// Import the package
import pubsubclient

// Publisher config requires a project name.
publisherConfig := pubsubclient.PublisherConfig{
    ProjectID: projectName,
}

// Passing the context and a config to GetPublisher fetches a publisher and errors
ctx := context.Background()
publisher, err := pubsubclient.GetPublisher(ctx, publisherConfig)

// Handle any errors as fatal and do not proceed further
if err != nil {
    log.Fatalf("Error occured while creating a publisher, Err: %v", err)
}

// All's good go ahead and publish

// StopAll method will be used to publish all the messages and stop all
// the goroutines that handle publish. There by gauranteeing the messages
// are flushed out from the internal queue to google pub/sub
defer publisher.StopAll()

for i := 0; i < 100; i++ {
    // Publish returns an error. The publish method also gives the ability
    // to publish the `same message` to more than `one topic`
    id, err := publisher.Publish(ctx, i, []string{topicName})
    if err != nil {
        log.Printf("Error occured while publishing the message, Err: %v", err)
    }
    log.Printf("MessageID: %v", id)
}
```

## Creating the subscriber

```go
// Imports the package
import pubsubclient

// Subscriber config requires a ProjectID, topicname, subscriptionName, ErrorHandler to
// handle errors, Handle func to process the messages

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

// Passing the context and config to CreateSubscription fetches a subscriber and errors
ctx := context.Background()
subscriber, err := pubsubclient.CreateSubscription(ctx, subscriberConfig)

// Handle any errors as fatal and do not proceed further
if err != nil {
    log.Fatalf("Error occured while creating a subscriber, Err: %v", err)
}

// Process will start pulling from the pubsub. The process accepts a waitgroup as
// it will be easier for us to orchestrate a use case where one application needs
// more than one subscriber
var wg sync.WaitGroup
wg.Add(1)
go subscriber.Process(ctx, &wg)
publishMessages(ctx, publisher)
publishMessages(ctx, publisher)
wg.Wait()
```

For a good idea on how to use the client please check the `pubsub_test.go` file.
