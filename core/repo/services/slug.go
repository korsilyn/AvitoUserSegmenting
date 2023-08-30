package service

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/repo"
	"context"
)

type SlugService struct {
	slugRepo repo.Slug
}

func NewSlugService(slugRepo repo.Slug) *SlugService {
	return &SlugService{slugRepo: slugRepo}
}

func (s *SlugService) CreateSlug(ctx context.Context, name string) (int, error) {
	var slug entity.Slug
	slug.Name = name
	return s.slugRepo.CreateSlug(ctx, slug)
}

func (s *SlugService) RemoveSlugByName(ctx context.Context, name string) error {
	return s.slugRepo.RemoveSlugByName(ctx, name)
}
