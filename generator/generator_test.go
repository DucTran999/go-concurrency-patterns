package generator_test

import (
	"context"
	"testing"
	"time"

	"github.com/DucTran999/go-concurrency-patterns/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GeneratorProducesValues(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := generator.Generator(ctx)

	var values []int
	for range 2 {
		select {
		case nonce, ok := <-ch:
			if !ok {
				t.Fatalf("channel closed prematurely")
			}

			assert.True(t, nonce >= 0 && nonce < 26)

			// Store nonce
			values = append(values, nonce)
		case <-time.After(3 * time.Second):
			t.Fatalf("timed out waiting for generator output")
		}
	}

	require.NotEmpty(t, values)
}
