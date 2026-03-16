package command

import (
	"BaseProjectGolang/internal/config"
	internallog "BaseProjectGolang/pkg/log"
	"context"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/soner3/flora"
)

type Scheduler struct {
	flora.Component
	cfg          *config.Config
	CronInstance *cron.Cron
	cancelFunc   context.CancelFunc
	mu           sync.Mutex
	running      bool
	logger       *internallog.Logger
}

func NewScheduler(
	cfg *config.Config,
	logger *internallog.Logger,
) *Scheduler {
	cronInstance := cron.New(cron.WithLogger(
		cron.VerbosePrintfLogger(logger.Logger),
	))

	return &Scheduler{
		cfg:          cfg,
		CronInstance: cronInstance,
		logger:       logger,
	}
}

// todo Реализовать интерфейс и вынести в отдельный пакет для переопределния в случае необходимости
func (scheduler *Scheduler) Schedule(ctx context.Context) {
	scheduler.mu.Lock()

	if scheduler.running {
		scheduler.mu.Unlock()
		return
	}

	scheduler.running = true

	// Example cancellable context
	ctx, scheduler.cancelFunc = context.WithCancel(ctx)
	scheduler.mu.Unlock()

	log.Println("Starting Cron jobs...")
	//Add schedule methods

	// Start the cron scheduler
	scheduler.CronInstance.Start()

	// Wait for context cancellation
	<-ctx.Done()

	// Stop cron when context is cancelled
	scheduler.CronInstance.Stop()
	scheduler.mu.Lock()
	scheduler.running = false
	scheduler.mu.Unlock()
	log.Println("Scheduler stopped")
}

func (scheduler *Scheduler) Stop() {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()

	if scheduler.cancelFunc != nil {
		scheduler.cancelFunc() // This will unblock the Schedule method
	}
}
