package events

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strconv"
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
	log.Info().Str("maxConcurrentTasks", strconv.Itoa(options.MaxConcurrentTasks)).Msg("Event Pool created successfully")

	return &pool
}

func runTask(task func()) {
	taskId := uuid.New()
	base64ID := base64.RawURLEncoding.EncodeToString([]byte(taskId.String()))
	log.Debug().Str("task_id", base64ID).Msgf("Initiating an incoming task.")
	task()
	log.Debug().Str("task_id", base64ID).Msgf("Task finished running.")
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
				runTask(snapshot.events[i])
				r.activeTaskCount -= 1
				r.Exec()
			}()
		}

	}()
}

func (r *ConcurrentPool[E]) Add(event E) {

	if r.activeTaskCount >= r.config.MaxConcurrentTasks {
		if len(r.events) > 0 {
			r.events = append(r.events, event)
		} else {
			r.events = append(make([]E, 0), event)
		}
		return
	}

	r.activeTaskCount += 1
	go func() {
		runTask(event)
		r.activeTaskCount -= 1
		r.Exec()
	}()
}
