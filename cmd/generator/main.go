package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucTran999/go-concurrency-patterns/generator"
)

// This app used generator pattern for generating
// the nonce a random number used for encrypt and decrypt message
// in Bob and Alice conversation
func main() {
	bobSentChan := make(chan string)
	aliceSentChan := make(chan string)

	bob := generator.NewPerson("Bob", []string{
		"Hi, Alice",
		"Would you like to join me for dinner at the restaurant tonight at 7:00 pm?",
	}, bobSentChan, aliceSentChan, 0)
	alice := generator.NewPerson("Alice", []string{
		"Hi, Bob",
		"Sounds great! See you at 7!",
	}, aliceSentChan, bobSentChan, 1)

	// Main context
	ctx, cancel := context.WithCancel(context.Background())

	secretChan := generator.Generator(ctx)
	go func() {
		for s := range secretChan {
			secret := s
			bob.SecretChan <- secret
			alice.SecretChan <- secret
		}
	}()

	go bob.Chat()
	go alice.Chat()

	closeApp(cancel)
}

func closeApp(cancel context.CancelFunc) {
	// Create a context that listens for interrupt or terminate signals
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // Clean up signal handlers when function exits

	<-shutdownCtx.Done() // Block until a signal is received
	log.Println("Shutdown signal received...")

	// Cancel the main context (e.g., passed to workers, generators, etc.)
	cancel()

	// Wait a bit to allow graceful cleanup
	time.Sleep(2 * time.Second)
	log.Println("shutdown app")
}
