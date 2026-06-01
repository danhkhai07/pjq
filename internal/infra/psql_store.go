package infra

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"pjq/internal/domain"

	_ "github.com/lib/pq"
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

func (s *PSQLStore) Save(job domain.Job) error {
	logs, _ := json.Marshal(job.Logs)
	_, err := s.db.Exec(`
		INSERT INTO jobs(id, type, payload, status, priority, retries, max_retries, error, result, logs, created_at, started_at, finished_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13)
		ON CONFLICT(id) DO UPDATE SET
			status      = EXCLUDED.status,
            retries     = EXCLUDED.retries,
            error       = EXCLUDED.error,
            result      = EXCLUDED.result,
            logs        = EXCLUDED.logs,
            started_at  = EXCLUDED.started_at,
            finished_at = EXCLUDED.finished_at
	`,
		job.ID, job.Type, job.Payload, job.Status,
		job.Priority, job.Retries, job.MaxRetries,
        job.Error, job.Result, logs, job.CreatedAt,
		job.StartedAt, job.FinishedAt,
	)
	return err
}

func (s *PSQLStore) Get(id string) (domain.Job, error) {
	row := s.db.QueryRow(selectStarJob() + ` WHERE id=$1;`, id)
	return scanJob(row)
}

func (s *PSQLStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM jobs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PSQLStore) List(filter domain.JobFilter) ([]domain.Job, error) {
	query := selectStarJob() + ` WHERE 1=1`
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

	rows, err := s.db.Query(query, args...)
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

func selectStarJob() string {
	return `SELECT id, type, payload, status, priority, retries, max_retries, error, result, logs, created_at, started_at, finished_at FROM jobs`

}

func scanJob(row interface{ Scan (...any) error }) (domain.Job, error) {
	var job domain.Job
	var logs []byte
	err := row.Scan(
		&job.ID, &job.Type, &job.Payload, &job.Status,
        &job.Priority, &job.Retries, &job.MaxRetries,
        &job.Error, &job.Result, &logs,
        &job.CreatedAt, &job.StartedAt, &job.FinishedAt,
	)
	if err != nil {
		return domain.Job{}, err
	}
	json.Unmarshal(logs, &job.Logs)
	return job, nil
}
