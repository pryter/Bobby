package events

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strconv"
)

// ConcurrentPoolOptions is options for ConcurrentPool
type ConcurrentPoolOptions struct {
	MaxConcurrentTasks int
}

// ConcurrentPool manages concurrent tasks by queueing tasks that exceed
// the MaxConcurrentTasks limit.
// This can be initiated using InitConcurrentPool method and should not
// be created by filling the struct.
type ConcurrentPool[E func()] struct {
	events          []E
	config          ConcurrentPoolOptions
	activeTaskCount int
	keyFunc         func(p any) string
}

// InitConcurrentPool creates ConcurrentPool
func InitConcurrentPool[E func()](options ConcurrentPoolOptions) *ConcurrentPool[E] {
	pool := ConcurrentPool[E]{config: options, activeTaskCount: 0}
	log.Info().Str(
		"maxConcurrentTasks", strconv.Itoa(options.MaxConcurrentTasks),
	).Msg("Event Pool created successfully")

	return &pool
}

// runTask runs provided task and track the task processes using uuid.
func runTask(task func()) {
	taskId := uuid.New()
	base64ID := base64.RawURLEncoding.EncodeToString([]byte(taskId.String()))
	log.Debug().Str("task_id", base64ID).Msgf("Initiating an incoming task.")
	task()
	log.Debug().Str("task_id", base64ID).Msgf("Task finished running.")
}

// Exec method executes tasks in the queue according to the available space and
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

// Add method adds a given event to the queue or if running tasks are not
// exceeded the limit it will run the task immediately.
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
