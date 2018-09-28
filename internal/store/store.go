package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AddPool(ctx context.Context) (int, error) {
	var id int
	err := s.db.GetContext(ctx, &id, "insert into pools default values returning id")
	return id, errors.WithStack(err)
}

func (s *Store) ListPools(ctx context.Context) ([]int, error) {
	var pools []int
	err := s.db.SelectContext(ctx, &pools, "select id from pools order by id")
	return pools, errors.WithStack(err)
}

type Job struct {
	Count   int
	Name    string
	PoolIDs string `db:"pool_ids"`
}

func (s *Store) ListJobs(ctx context.Context) ([]Job, error) {
	var jobs []Job
	err := s.db.SelectContext(ctx, &jobs, "select count(*), name, string_agg(pool_id::text, ', ') as pool_ids from jobs group by name")
	return jobs, errors.WithStack(err)
}

func (s *Store) AddJob(ctx context.Context, name string, poolID int) (int, error) {
	var id int
	err := s.db.GetContext(ctx, &id, "insert into jobs(name, pool_id) values($1, $2) returning id", name, poolID)
	return id, errors.WithStack(err)
}
