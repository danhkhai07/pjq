package infra

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"context"

	"pjq/internal/domain"

	_ "github.com/lib/pq"
)

const (
	SELECT_STAR_JOB = `SELECT id, type, payload, status, priority, retries, max_retries, error, result, logs, created_at, started_at, finished_at, run_at FROM jobs `
)

type PSQLStore struct {
	db *sql.DB
}

func NewPSQLStore(connStr string) (*PSQLStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PSQLStore{ db }, nil
}

func (s *PSQLStore) Save(ctx context.Context, job domain.Job) error {
	logs, _ := json.Marshal(job.Logs)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO jobs(id, type, payload, status, priority, retries, max_retries, error, result, logs, created_at, started_at, finished_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13)
		ON CONFLICT(id) DO UPDATE SET
			status      = EXCLUDED.status,
            retries     = EXCLUDED.retries,
            error       = EXCLUDED.error,
            result      = EXCLUDED.result,
            logs        = EXCLUDED.logs,
            started_at  = EXCLUDED.started_at,
            finished_at = EXCLUDED.finished_at,
			run_at		= EXCLUDED.run_at
	`,
		job.ID, job.Type, job.Payload, job.Status,
		job.Priority, job.Retries, job.MaxRetries,
        job.Error, job.Result, logs, job.CreatedAt,
		job.StartedAt, job.FinishedAt,
	)
	return err
}

func (s *PSQLStore) Get(ctx context.Context, id string) (domain.Job, error) {
	row := s.db.QueryRowContext(ctx, SELECT_STAR_JOB + ` WHERE id=$1;`, id)
	return scanJob(row)
}

func (s *PSQLStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PSQLStore) List(ctx context.Context, filter domain.JobFilter) ([]domain.Job, error) {
	query := SELECT_STAR_JOB + ` WHERE 1=1`
	args := []any{}
	i := 1

	if filter.Status != nil {
		query += fmt.Sprintf(` AND status = $%d`, i)
		args = append(args, *filter.Status)
		i++
	}
	if filter.Type != nil {
		query += fmt.Sprintf(` AND type = $%d`, i)
		args = append(args, *filter.Type)
		i++
	}
	if filter.Status != nil {
		query += ` AND retries < max_retries`
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *PSQLStore) Recover(ctx context.Context) ([]domain.Job, error) {
	rows, err := s.db.QueryContext(ctx,
		SELECT_STAR_JOB +
		` WHERE status in ('pending', 'running')`,
	)
	if err != nil {
		return nil, err
	}

	jobs := []domain.Job{}

	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func scanJob(row interface{ Scan (...any) error }) (domain.Job, error) {
	var job domain.Job
	var logs []byte
	err := row.Scan(
		&job.ID, &job.Type, &job.Payload, &job.Status,
        &job.Priority, &job.Retries, &job.MaxRetries,
        &job.Error, &job.Result, &logs,
        &job.CreatedAt, &job.StartedAt, &job.FinishedAt, &job.RunAt,
	)
	if err != nil {
		return domain.Job{}, err
	}
	json.Unmarshal(logs, &job.Logs)
	return job, nil
}
