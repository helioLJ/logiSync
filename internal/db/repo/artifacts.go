package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtifactRepo struct {
	pool *pgxpool.Pool
}

func NewArtifactRepo(pool *pgxpool.Pool) *ArtifactRepo {
	return &ArtifactRepo{pool: pool}
}

func (r *ArtifactRepo) Insert(ctx context.Context, jobID uuid.UUID, provider, key, kind string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO artifacts (id, job_id, provider, artifact_key, kind)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
	`, jobID, provider, key, kind)
	if err != nil {
		return fmt.Errorf("insert artifact: %w", err)
	}
	return nil
}
