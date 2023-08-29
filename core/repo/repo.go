package repo

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/repo/pgdb"
	"avito-user-segmenting/modules/postgresql"
	"context"
)

type Slug interface {
	CreateSlug(ctx context.Context, slug entity.Slug) (int, error)
	RemoveSlugByName(ctx context.Context, name string) (error)
	GetSlugById(ctx context.Context, id int) (entity.Slug, error)
	GetSlugByName(ctx context.Context, name string) (entity.Slug, error)
}

type Operation interface {
	CreateOperation(ctx context.Context, operation entity.Operation) (int, error)
	RemoveOperationByUserId(ctx context.Context, userId int, slugId int) (error)
	GetOperationById(ctx context.Context, id int) (entity.Operation, error)
	GetAllSlugsByUserId(ctx context.Context, userId int) ([]int, error)
}

type Repositories struct {
	Slug
	Operation
}

func NewRepositories(pg *postgresql.PostgreSQL) *Repositories {
	return &Repositories{
		Slug: pgdb.NewSlugRepo(pg),
		Operation: pgdb.NewOperationRepo(pg),
	}
}
