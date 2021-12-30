// Package pubsub knows how to publish and subscribe messages to/from a broker
package pubsub

// Publisher knows how to publish messages
type Publisher interface {
	// SendMsg sends messages
	SendMsg(msg string) error

	// Close closes the publisher
	Close() error
}
