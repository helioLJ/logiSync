package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResultRepo struct {
	pool *pgxpool.Pool
}

type TrackingResult struct {
	Provider          string          `json:"provider"`
	TrackingCode      string          `json:"tracking_code"`
	NormalizedPayload json.RawMessage `json:"normalized_payload"`
	CreatedAt         string          `json:"created_at"`
}

func NewResultRepo(pool *pgxpool.Pool) *ResultRepo {
	return &ResultRepo{pool: pool}
}

func (r *ResultRepo) Insert(ctx context.Context, jobID uuid.UUID, provider, trackingCode string, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO tracking_results (id, job_id, provider, tracking_code, normalized_payload)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
	`, jobID, provider, trackingCode, payloadJSON)
	if err != nil {
		return fmt.Errorf("insert tracking result: %w", err)
	}
	return nil
}

func (r *ResultRepo) GetLatestByJobID(ctx context.Context, jobID uuid.UUID) (TrackingResult, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT provider, tracking_code, normalized_payload, created_at
		FROM tracking_results
		WHERE job_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, jobID)

	var result TrackingResult
	var createdAt time.Time
	if err := row.Scan(&result.Provider, &result.TrackingCode, &result.NormalizedPayload, &createdAt); err != nil {
		if err == pgx.ErrNoRows {
			return TrackingResult{}, err
		}
		return TrackingResult{}, fmt.Errorf("get latest result: %w", err)
	}
	result.CreatedAt = createdAt.Format(time.RFC3339)
	return result, nil
}
