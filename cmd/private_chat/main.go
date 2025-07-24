package main

import (
	"context"
	"sync"

	privatechat "github.com/DucTran999/go-concurrency-patterns/case_study/private_chat"
)

// This app used generator pattern for generating
// the nonce a random number used for encrypt and decrypt message
// in Bob and Alice conversation
func main() {
	bobSentChan := make(chan string)
	aliceSentChan := make(chan string)

	bob := privatechat.NewPerson("Bob", []string{
		"Hi, Alice",
		"Would you like to join me for dinner at the restaurant tonight at 7:00 pm?",
	}, bobSentChan, aliceSentChan, 0)
	alice := privatechat.NewPerson("Alice", []string{
		"Hi, Bob",
		"Sounds great! See you at 7!",
	}, aliceSentChan, bobSentChan, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	secretChan := privatechat.Generator(ctx)
	go func() {
		for s := range secretChan {
			secret := s
			bob.SecretChan <- secret
			alice.SecretChan <- secret
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bob.Chat()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		alice.Chat()
	}()

	wg.Wait()
}
