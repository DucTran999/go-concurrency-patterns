package generator_test

import (
	"testing"

	"github.com/DucTran999/go-concurrency-patterns/generator"
)

func TestPersonChat(t *testing.T) {
	t.Parallel()

	senderToReceiver := make(chan string, 1)
	receiverToSender := make(chan string, 1)
	nonceChan := make(chan int, 10)
	done := make(chan bool, 2)

	sender := generator.NewPerson("Alice", []string{"Hello"}, senderToReceiver, receiverToSender, 0)
	receiver := generator.NewPerson("Bob", []string{"Hello"}, receiverToSender, senderToReceiver, 1)

	// Use a simple Caesar shift for testing
	// Start chat goroutines with completion signaling
	go func() {
		sender.Chat()
		done <- true
	}()

	go func() {
		receiver.Chat()
		done <- true
	}()

	// Provide a nonce and close the channel to end chat
	nonces := []int{3, 4, 5}
	for _, n := range nonces {
		receiver.SecretChan <- n
		sender.SecretChan <- n
	}

	// Close secret channels to terminate chat
	close(sender.SecretChan)
	close(receiver.SecretChan)
	close(nonceChan)

	// Wait for both goroutines to complete
	<-done
	<-done
}
