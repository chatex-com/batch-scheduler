package batch

import (
	"sync"
	"time"
)

type Scheduler struct {
	lock     sync.RWMutex
	schedule Schedule
	ch       chan interface{}
}

func NewScheduler(cfg []Rule) (*Scheduler, error) {
	schedule, err := NewSchedule(cfg)
	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		schedule: *schedule,
	}

	return scheduler, nil
}

func (s *Scheduler) Run() <-chan interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ch != nil {
		return s.ch
	}

	// Give the channel a 1-element buffer.
	// If the client falls behind while reading, we drop next ticks
	// on the floor until the client catches up.
	s.ch = make(chan interface{}, 1)

	go s.run()

	return s.ch
}

func (s *Scheduler) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	close(s.ch)
	s.ch = nil
}

func (s *Scheduler) run() {
	for {
		s.lock.Lock()
		s.ch <- true
		s.lock.Unlock()

		<-time.After(s.schedule.GetWindow(time.Now()))
	}
}
