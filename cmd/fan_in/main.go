package main

import (
	"log"

	fanin "github.com/DucTran999/go-concurrency-patterns/fan_in"
)

func main() {
	serviceA := make(chan string)
	serviceB := make(chan string)
	serviceC := make(chan string)

	// Start 3 log producers
	go fanin.LogProducer("ServiceA", serviceA)
	go fanin.LogProducer("ServiceB", serviceB)
	go fanin.LogProducer("ServiceC", serviceC)

	// Fan-in their output
	combined := fanin.MergeChanel(serviceA, serviceB, serviceC)

	// Read 9 log entries (3 per service)
	for range 9 {
		srvLog := <-combined
		log.Println("[INFO]", srvLog)
	}
}
