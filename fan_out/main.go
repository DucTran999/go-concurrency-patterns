package main

import (
	"log"
	"sync"
)

func worker(id int, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {
		log.Printf("[INFO] Worker %d received: %s\n", id, msg)
	}
}

func main() {
	const workerCount = 3
	message := "Hello from sender"

	var wg sync.WaitGroup
	channels := make([]chan string, workerCount)

	for i := range workerCount {
		ch := make(chan string)
		channels[i] = ch
		wg.Add(1)
		go worker(i+1, ch, &wg)
	}

	// Fan-out: send message to all channels (worker)
	for _, ch := range channels {
		ch <- message
		close(ch)
	}

	wg.Wait()
	log.Println("All workers finished.")
}
