package scheduler

import (
	"log"
	"time"

	"go-reminder/internal/service"
)

type Scheduler struct {
	reminderService service.ReminderService
	interval        time.Duration
	stop            chan struct{}
}

func New(reminderService service.ReminderService, interval time.Duration) *Scheduler {
	return &Scheduler{
		reminderService: reminderService,
		interval:        interval,
		stop:            make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		s.runOnce()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runOnce()
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) runOnce() {
	if err := s.reminderService.ProcessDueReminders(); err != nil {
		log.Printf("scheduler failed: %v", err)
	}
}
