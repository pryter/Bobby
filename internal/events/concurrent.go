package events

import (
	"fmt"
)

type ConcurrentPoolOptions struct {
	MaxConcurrentTasks int
}

type ConcurrentPool[E func()] struct {
	events          []E
	config          ConcurrentPoolOptions
	activeTaskCount int
	keyFunc         func(p any) string
}

func InitConcurrentPool[E func()](options ConcurrentPoolOptions) *ConcurrentPool[E] {
	pool := ConcurrentPool[E]{config: options, activeTaskCount: 0}

	return &pool
}

func (r *ConcurrentPool[E]) Exec() {

	eCount := len(r.events)
	slot := r.config.MaxConcurrentTasks - r.activeTaskCount

	if slot <= 0 {
		return
	}

	tasks := slot
	if tasks > eCount {
		tasks = eCount
	}

	r.activeTaskCount += tasks
	snapshot := *r
	r.events = append(r.events[:0], r.events[tasks:]...)

	go func() {
		for i := 0; i < tasks; i++ {
			i := i
			go func() {
				fmt.Printf("running task %d from queue\n", i)
				snapshot.events[i]()
				r.activeTaskCount -= 1
				r.Exec()
			}()
		}

	}()
}

func (r *ConcurrentPool[E]) Add(event E) {

	if r.activeTaskCount >= r.config.MaxConcurrentTasks {
		println("saved task to queue")
		if len(r.events) > 0 {
			r.events = append(r.events, event)
		} else {
			r.events = append(make([]E, 0), event)
		}
		return
	}

	r.activeTaskCount += 1
	go func() {
		fmt.Printf("running requested task \n")
		event()
		r.activeTaskCount -= 1
		r.Exec()
	}()
}
