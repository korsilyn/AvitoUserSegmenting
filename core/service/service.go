package service

import (
	"avito-user-segmenting/core/repo"
	"context"
)

type Slug interface {
	CreateSlug(ctx context.Context, name string) (int, error)
	RemoveSlugByName(ctx context.Context, name string) error
}

type OperationCreateInput struct {
	Slugs []string
	Percent int
	TTL int
}

type OperationRemoveInput struct {
	Slugs []string
	UserId int
}

type Operation interface {
	CreateOperations(ctx context.Context, input OperationCreateInput) error
	RemoveOperationsByUserId(ctx context.Context, input OperationRemoveInput) error
	GetAllSlugsByUserId(ctx context.Context, userId int) ([]string, error)
}

type Services struct {
	Operation Operation
	Slug Slug
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Operation: NewOperationService(deps.Repos.Operation, deps.Repos.Slug),
		Slug: NewSlugService(deps.Repos.Slug),
	}
}
