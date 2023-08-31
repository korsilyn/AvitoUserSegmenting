package pgdb

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/repo/repoerrors"
	"avito-user-segmenting/modules/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type OperationRepo struct {
	*postgresql.PostgreSQL
}

func NewOperationRepo(pg *postgresql.PostgreSQL) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) CreateOperation(ctx context.Context, operation entity.Operation) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("OperationRepo.CreateOperation - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	/*sql, args, _ := r.Builder.
		Select("id").
		From("operations").
		Where("user_id = ? and slug_id = ?", operation.UserId, operation.SlugId).
		ToSql()
	var tmp int
	err = tx.QueryRow(ctx, sql, args...).Scan(&tmp)
	if tmp != 0 {
		return 0, repoerrors.ErrAlreadyExists
	}*/

	sql, args, _ := r.Builder.
		Insert("operations").
		Columns("user_id", "slug_id", "created_at", "removed_at").
		Values(operation.UserId, operation.SlugId, time.Now(), operation.RemovedAt).
		Suffix("RETURNING id").
		ToSql()
	
	var id int
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("OperationRepo.CreateOperation - tx.QueryRow - %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("OperationRepo.CreateOperation - tx.Commit - %v", err)
	}

	return id, nil
}

func (r *OperationRepo) RemoveOperationByUserId(ctx context.Context, userId int, slugId int) error {
	sql, args, _ := r.Builder.
		Update("operations").
		Where("user_id = ? and slug_id = ?", userId, slugId).
		Set("removed_at", time.Now()).
		ToSql()
	
	_, err := r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return repoerrors.ErrNotFound
	}

	return nil
}

func (r *OperationRepo) GetOperationById(ctx context.Context, id int) (entity.Operation, error) {
	sql, args, _ := r.Builder.
		Select("id, user_id, slug_id, added_at, removed_at").
		From("operations").
		Where("id = ? and (removed_at > ? or removed_at = ?)", id, time.Now(), time.Time{}).
		ToSql()

	var operation entity.Operation
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&operation.Id,
		&operation.UserId,
		&operation.SlugId,
		&operation.AddedAt,
		&operation.RemovedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Operation{}, repoerrors.ErrNotFound
		}
		return entity.Operation{}, fmt.Errorf("OperationRepo.GetOperationById - r.Pool.QueryRow: %v", err)
	}
	
	return operation, nil
}

func (r *OperationRepo) GetAllSlugsByUserId(ctx context.Context, userId int) ([]int, error) {
	sql, args, _ := r.Builder.
		Select("slug_id").
		From("operations").
		Where("user_id = ? and removed_at = ?", userId, time.Time{}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("OperationRepo.GetAllSlugsByUserId - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var slugs []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("OperationRepo.GetAllSlugsByUserId - rows.Scan: %v", err)
		}

		slugs = append(slugs, id)
	}

	return slugs, nil
}
