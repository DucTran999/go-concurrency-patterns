package fanin

import (
	"fmt"
	"sync"
	"time"
)

// Fan-in: Merge multiple input channels into one output channel
func MergeChanel(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	merged := make(chan string)

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan string) {
			defer wg.Done()
			for msg := range ch {
				merged <- msg
			}
		}(ch) // Pass loop var into closure
	}

	// Wait for all senders to finish, then close the merged channel
	go func() {
		wg.Wait()
		close(merged)
	}()

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
