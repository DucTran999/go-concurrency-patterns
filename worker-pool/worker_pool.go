package workerpool

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type workerPool struct {
	numOfWorkers int

	jobs    chan *Job
	results chan int64
	wg      sync.WaitGroup
}

func NewWorkerPool(numOfWorker int) *workerPool {
	return &workerPool{
		numOfWorkers: numOfWorker,
		jobs:         make(chan *Job, 100),
		results:      make(chan int64, 100),
		wg:           sync.WaitGroup{},
	}
}

func (wp *workerPool) StreamJobFromFile(filename string) error {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Println("failed to close file:", err)
		}
	}()

	reader := csv.NewReader(file)

	var lineNum int64 = 0
	for {
		record, err := reader.Read()
		if err != nil {
			// EOF will break the loop
			if errors.Is(err, io.EOF) {
				break
			}

			// Skip bad row, log and continue
			log.Printf("[WARN] line %d: %v", lineNum, err)
			lineNum++
			continue
		}

		// Process valid record
		wp.jobs <- NewJob(lineNum, record)
		lineNum++
	}
	close(wp.jobs)

	return nil
}

func (wp *workerPool) SpawnWorkers() {
	for id := range wp.numOfWorkers {
		wp.wg.Add(1)

		go wp.spawnWorker(id)
	}
}

func (wp *workerPool) CollectResult() {
	// Close results channel when all workers are done
	done := make(chan struct{})

	go func() {
		defer close(done)
		for jobErrId := range wp.results {
			log.Println("[ERROR] failed when process line:", jobErrId)
		}
	}()

	wp.wg.Wait()
	close(wp.results)
	<-done // Wait for all results to be processed
}

func (wp *workerPool) spawnWorker(workerID int) {
	log.Printf("[INFO] worker %d running...", workerID)
	defer wp.wg.Done()

	for job := range wp.jobs {
		err := job.Run()

		if err != nil {
			wp.results <- job.ID
		}
	}
}
