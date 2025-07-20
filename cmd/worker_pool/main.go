package main

import (
	"log"
	"runtime"
	"time"

	"github.com/DucTran999/go-concurrency-patterns/test/testutil"
	workerpool "github.com/DucTran999/go-concurrency-patterns/worker-pool"
)

func main() {

	s := time.Now()

	path, err := testutil.BuildFilePath("test/datatest/mock_data.csv")
	if err != nil {
		log.Fatalln("failed to lookup file path", err)
	}

	// Initialized worker pool
	numberOfWorker := runtime.NumCPU() * 2
	wp := workerpool.NewWorkerPool(numberOfWorker)

	go func() {
		if err = wp.StreamJobFromFile(path); err != nil {
			log.Fatalln("failed to stream file", err)
		}
	}()

	// spawn worker process line record
	wp.SpawnWorkers()

	// Collect line error from worker and wait for all worker done job
	wp.CollectResult()

	log.Println("Processed time:", time.Since(s))
}
