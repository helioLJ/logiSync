package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Job struct {
	ID           uuid.UUID `json:"job_id"`
	Provider     string    `json:"provider"`
	TrackingCode string    `json:"tracking_code"`
	Status       string    `json:"status"`
	Attempts     int       `json:"attempts"`
	ErrorCode    *string   `json:"error_code,omitempty"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type JobRepo struct {
	pool *pgxpool.Pool
}

func NewJobRepo(pool *pgxpool.Pool) *JobRepo {
	return &JobRepo{pool: pool}
}

func (r *JobRepo) Create(ctx context.Context, job Job) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO jobs (id, provider, tracking_code, status, attempts)
		VALUES ($1, $2, $3, $4, $5)
	`, job.ID, job.Provider, job.TrackingCode, job.Status, job.Attempts)
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}
	return nil
}

func (r *JobRepo) Get(ctx context.Context, id uuid.UUID) (Job, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, provider, tracking_code, status, attempts, error_code, error_message, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`, id)

	var job Job
	if err := row.Scan(
		&job.ID,
		&job.Provider,
		&job.TrackingCode,
		&job.Status,
		&job.Attempts,
		&job.ErrorCode,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
	); err != nil {
		return Job{}, err
	}

	return job, nil
}

func (r *JobRepo) MarkRunning(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE jobs
		SET status = 'RUNNING', attempts = attempts + 1, updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark running: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *JobRepo) MarkDone(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE jobs
		SET status = 'DONE', updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark done: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *JobRepo) MarkFailed(ctx context.Context, id uuid.UUID, code, message string) error {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE jobs
		SET status = 'FAILED', error_code = $2, error_message = $3, updated_at = now()
		WHERE id = $1
	`, id, code, message)
	if err != nil {
		return fmt.Errorf("mark failed: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
