package main

import (
	"context"
	"log"
	"sync"
	"time"

	privatechat "github.com/DucTran999/go-concurrency-patterns/case_study/private_chat"
)

func main() {
	// Communication channels
	bobSentChan := make(chan string)
	aliceSentChan := make(chan string)

	// Create Bob and Alice with their messages
	bob := privatechat.NewPerson(
		"Bob",
		[]string{
			"Hi, Alice",
			"Would you like to join me for dinner at the restaurant tonight at 7:00 pm?",
		},
		bobSentChan, aliceSentChan, 0,
	)

	alice := privatechat.NewPerson(
		"Alice",
		[]string{
			"Hi, Bob",
			"Sounds great! See you at 7!",
		},
		aliceSentChan, bobSentChan, 1,
	)

	// Secret nonce generator
	ctx, cancel := context.WithCancel(context.Background())

	secretChan := privatechat.Generator(ctx)

	// Fan-out the generated secrets to both Alice and Bob
	go func() {
		for secret := range secretChan {
			bob.ReceiveSecret(secret)
			alice.ReceiveSecret(secret)
		}
	}()

	// Start chatting
	log.Println("[INFO] conversation start")
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bob.Chat()
	}()

	go func() {
		defer wg.Done()
		alice.Chat()
	}()

	wg.Wait()

	cancel()                           // signal to stop secret generator if it's still running
	time.Sleep(100 * time.Millisecond) // Give time for graceful closing
	log.Println("[INFO] conversation ended.")
}
