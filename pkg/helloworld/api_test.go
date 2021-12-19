package helloworld_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yngvark/go-rest-api-template/pkg/helloworld"
)

func TestHello(t *testing.T) {
	testCases := []struct {
		name   string
		expect string
	}{
		{
			name:   "Should return hello world",
			expect: "Hello world!",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, helloworld.Hello())
		})
	}
}
