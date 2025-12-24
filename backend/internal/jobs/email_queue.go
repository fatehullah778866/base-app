package jobs

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"base-app-service/internal/services"
)

type EmailJob struct {
	To      string
	Subject string
	Body    string
	HTML    string
}

type JobQueue struct {
	emailService *services.EmailService
	logger       *zap.Logger
	jobs         chan EmailJob
	workers      int
}

func NewJobQueue(emailService *services.EmailService, logger *zap.Logger, workers int) *JobQueue {
	if workers <= 0 {
		workers = 3 // Default 3 workers
	}

	queue := &JobQueue{
		emailService: emailService,
		logger:       logger,
		jobs:         make(chan EmailJob, 100), // Buffer 100 jobs
		workers:      workers,
	}

	// Start workers
	for i := 0; i < workers; i++ {
		go queue.worker(i)
	}

	return queue
}

func (q *JobQueue) EnqueueEmail(ctx context.Context, job EmailJob) error {
	select {
	case q.jobs <- job:
		q.logger.Info("Email job enqueued", zap.String("to", job.To))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("job queue full")
	}
}

func (q *JobQueue) worker(id int) {
	for job := range q.jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		
		email := services.Email{
			To:      []string{job.To},
			Subject: job.Subject,
			HTML:    job.HTML,
			Body:    job.Body,
		}

		if err := q.emailService.SendEmail(ctx, email); err != nil {
			q.logger.Error("Failed to process email job",
				zap.Int("worker_id", id),
				zap.String("to", job.To),
				zap.Error(err),
			)
		} else {
			q.logger.Info("Email job processed successfully",
				zap.Int("worker_id", id),
				zap.String("to", job.To),
			)
		}

		cancel()
	}
}

func (q *JobQueue) Shutdown() {
	close(q.jobs)
}

