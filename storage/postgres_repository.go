package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"shortenLink/models"
)

type PostgresRepository struct {
	db *sqlx.DB
}

// pkg/storage/postgres_repository.go
func (r *PostgresRepository) CreateShortURL(ctx context.Context, su *models.ShortUrl) error {
	const query = `
        INSERT INTO short_urls
            (short_code, origin_url, expires_at)
        VALUES
            (:short_code, :origin_url, :expires_at)
        ON CONFLICT (short_code)
        DO UPDATE SET
            origin_url = EXCLUDED.origin_url,
            expires_at = EXCLUDED.expires_at,
            deleted_at = NULL
        RETURNING id, created_at`

	_, err := r.db.NamedExecContext(ctx, query, su)
	return err
}

func (r *PostgresRepository) GetOriginalURL(ctx context.Context, code string) (string, error) {
	const query = `
        SELECT origin_url
        FROM short_urls
        WHERE short_code = $1
          AND deleted_at IS NULL
          AND (expires_at IS NULL OR expires_at > NOW())`

	var url string
	err := r.db.GetContext(ctx, &url, query, code)

	return url, err
}

func (r *PostgresRepository) Save(url *models.ShortUrl) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.NamedExec(`
        INSERT INTO short_urls
        (short_code, origin_url, expires_at)
        VALUES (:short_code, :origin_url, :expires_at)
    `, url); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) UpdateStats(code string) error {
	return r.db.Get(&models.VisitStats{}, `
        INSERT INTO visit_stats (short_code, visit_count)
        VALUES ($1, 1)
        ON CONFLICT (short_code)
        DO UPDATE SET
            visit_count = visit_stats.visit_count + 1,
            last_visit = NOW()
    `, code)
}
