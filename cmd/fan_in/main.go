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
	go fanin.LogProducer("Service A", serviceA)
	go fanin.LogProducer("Service B", serviceB)
	go fanin.LogProducer("Service C", serviceC)

	// Fan-in their output
	combined := fanin.MergeChanel(serviceA, serviceB, serviceC)

	for msg := range combined {
		log.Println("[INFO]", msg)
	}
}
