package dht

import (
	"time"
)

type Scheduler struct {
	quit chan struct{}
}

type fn func()

func NewScheduler(k *Kademlia) *Scheduler {
	return &Scheduler{quit: make(chan struct{})}
}

func (s *Scheduler) RepeatTask(interval int, task fn) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			task()
		case <-s.quit:
			ticker.Stop()
			return
		}
	}
}
