package pubsub

// Consumer knows how to consume messages
type Consumer interface {
	// ListenForMessages starts to receive messages which will be available by reading SubscriberChannel(). It blocks
	// until the Consumer's context is canceled, so you should start it as a goroutine.
	ListenForMessages() error

	// SubscriberChannel returns a channel which can be used for reading incoming messages, or for sending messages to.
	SubscriberChannel() chan string

	// Close closes the Consumer
	Close() error
}
