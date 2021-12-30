package broadcast_test

import (
	"testing"

	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/broadcast"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Run("Should send message to listeners", func(t *testing.T) {
		// Given
		broadcaster := broadcast.New()
		testSubscriber := make(chan string)

		broadcaster.AddSubscriber(testSubscriber)

		// When
		go func() {
			err := broadcaster.BroadCast("hi")
			require.NoError(t, err)
		}()

		// Then
		lastMsgReceived := <-testSubscriber
		assert.Equal(t, "hi", lastMsgReceived)
	})
}
