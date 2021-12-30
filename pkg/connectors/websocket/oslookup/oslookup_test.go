package oslookup_test

import (
	"testing"

	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/log2"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/oslookup"

	"github.com/stretchr/testify/assert"
)

func TestOslookup(t *testing.T) {
	t.Run("Should parse cors worigins", func(t *testing.T) {
		logger, err := log2.New()
		assert.Nil(t, err)

		corsHelper := oslookup.NewCORSHelper(logger)
		allowed, err := corsHelper.GetAllowedCorsOrigins(osLookupEnv, "TEST_ENV")
		assert.Nil(t, err)

		expected := make(map[string]bool)
		expected["http://localhost:3000"] = true
		expected["https://localhost:3001"] = true

		assert.Equal(t, expected, allowed)
	})
}

func osLookupEnv(_ string) (string, bool) {
	return "http://localhost:3000,https://localhost:3001", true
}
