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

// Passing the config to GetPublisher fetches a publisher and errors
publisher, err := pubsubclient.GetPublisher(publisherConfig)

// Handle any errors as fatal and do not proceed further
if err != nil {
    log.Fatalf("Error occured while creating a publisher, Err: %v", err)
}

// All's good go ahead and publish
for i := 0; i < 100; i++ {
    // Publish returns a messageID and an error. MessageID is the message handle
    // returned by google pubsub. The publish method also gives the ability
    // to publish the `same message` to more than `one topic`
    id, err := publisher.Publish(i, []string{topicName})
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

// Passing the config to CreateSubscription fetches a subscriber and errors
subscriber, err := pubsubclient.CreateSubscription(subscriberConfig)

// Handle any errors as fatal and do not proceed further
if err != nil {
    log.Fatalf("Error occured while creating a subscriber, Err: %v", err)
}

// Process will start pulling from the pubsub. The process accepts a waitgroup as
// it will be easier for us to orchestrate a use case where one application needs
// more than one subscriber
var wg sync.WaitGroup
wg.Add(1)
go subscriber.Process(&wg)
publishMessages(publisher)
wg.Wait()
```

For a good idea on how to use the client please check the `pubsub_test.go` file.
