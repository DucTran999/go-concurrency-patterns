package fanin

import (
	"fmt"
	"time"
)

// Fan-in: Merge multiple input channels into one output channel
func MergeChanel(channels ...<-chan string) <-chan string {
	merged := make(chan string)

	for _, ch := range channels {
		go func(ch <-chan string) {
			for msg := range ch {
				merged <- msg
			}
		}(ch) // Pass loop var into closure
	}

	return merged
}

// Simulate a service that sends logs into its own channel
func LogProducer(serviceName string, out chan<- string) {
	for i := 1; i <= 3; i++ {
		out <- fmt.Sprintf("[%s] log %d", serviceName, i)
		time.Sleep(time.Millisecond * 300)
	}
	close(out)
}
