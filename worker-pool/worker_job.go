package workerpool

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Job struct {
	ID     int64
	record []string
}

func NewJob(id int64, record []string) *Job {
	return &Job{
		ID:     id,
		record: record,
	}
}

func (j *Job) Run() error {
	// Simulate process a record need 200ms
	time.Sleep(time.Millisecond * 200)

	if len(j.record) == 6 && !isIPv4(j.record[5]) {
		return fmt.Errorf("record %d: invalid IP", j.ID)
	}

	return nil
}

func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && strings.Count(ip, ":") == 0
}
