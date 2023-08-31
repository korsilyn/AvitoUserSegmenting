package pgdb

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/repo/repoerrors"
	"avito-user-segmenting/modules/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SlugRepo struct {
	*postgresql.PostgreSQL
}

func NewSlugRepo(pg *postgresql.PostgreSQL) *SlugRepo {
	return &SlugRepo{pg}
}

func (r *SlugRepo) CreateSlug(ctx context.Context, slug entity.Slug) (int, error) {
	sql, args, _ := r.Builder.
		Insert("slugs").
		Columns("name").
		Values(slug.Name).
		Suffix("RETURNING id").
		ToSql()
	
	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrors.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("SlugRepo.CreateSlug - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *SlugRepo) RemoveSlugByName(ctx context.Context, name string) error {
	sql, args, _ := r.Builder.
		Delete("slugs").
		Where("name = ?", name).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return repoerrors.ErrNotFound
		//return fmt.Errorf("SlugRepo.RemoveSlugByName - r.Pool.Exec: %v", err)
	}
	
	return nil
}

func (r *SlugRepo) GetSlugById(ctx context.Context, id int) (entity.Slug, error) {
	sql, args, _ := r.Builder.
		Select("id, name").
		From("slugs").
		Where("id = ?", id).
		ToSql()

	var slug entity.Slug
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&slug.Id,
		&slug.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Slug{}, repoerrors.ErrNotFound
		}
		return entity.Slug{}, fmt.Errorf("SlugRepo.GetSlugById - r.Pool.QueryRow: %v", err)
	}
	
	return slug, nil
}

func (r *SlugRepo) GetSlugByName(ctx context.Context, name string) (entity.Slug, error) {
	sql, args, _ := r.Builder.
		Select("id, name").
		From("slugs").
		Where("name = ?", name).
		ToSql()

	var slug entity.Slug
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&slug.Id,
		&slug.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Slug{}, repoerrors.ErrNotFound
		}
		return entity.Slug{}, fmt.Errorf("SlugRepo.GetSlugById - r.Pool.QueryRow: %v", err)
	}
	
	return slug, nil
}
