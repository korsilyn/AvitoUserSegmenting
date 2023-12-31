package service

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/repo"
	"context"
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type OperationService struct {
	operationRepo repo.Operation
	slugRepo repo.Slug
}

func NewOperationService(operationRepo repo.Operation, slugRepo repo.Slug) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
		slugRepo: slugRepo,
	}
}

func (s *OperationService)CreateOperations(ctx context.Context, input OperationCreateInput) error {
	for _, name := range input.Slugs {
		var operation entity.Operation
		var slug entity.Slug
		slug, err := s.slugRepo.GetSlugByName(ctx, name)
		if err != nil {
			return fmt.Errorf("OperationService.CreateOperations - GetSlugByName: %v", err)
		}
		operation.SlugId = slug.Id
		if input.TTL > 0 {
			now := time.Now()
			operation.RemovedAt = now.Add(time.Duration(input.TTL) * time.Hour)
		} else {
			operation.RemovedAt = time.Time{}
		}
		
		// Из-за того, что я взял пользователей 1000-1099
		ids := make([]int, 0, input.Percent)
		for i := 0; i < input.Percent; i++ {
			rand.Seed(time.Now().UnixNano())
			min := 1000
			max := 1099
			id := rand.Intn(max - min + 1) + min
			if slices.Contains(ids, id) {
				i--
				continue
			}
			ids = append(ids, id)
			operation.UserId = id
			_, err = s.operationRepo.CreateOperation(ctx, operation)
			if err != nil {
				return fmt.Errorf("OperationService.CreateOperation - CreateOperation: %v", err)
			}
		}
	}

	return nil
}

func (s *OperationService)RemoveOperationsByUserId(ctx context.Context, input OperationRemoveInput) error {
	for _, name := range input.Slugs {
		var slug entity.Slug
		slug, err := s.slugRepo.GetSlugByName(ctx, name)
		if err != nil {
			return fmt.Errorf("OperationService.RemoveOperationsByUserId - GetSlugByName: %v", err)
		}
		err = s.operationRepo.RemoveOperationByUserId(ctx, input.UserId, slug.Id)
		if err != nil {
			return fmt.Errorf("OperationService.RemoveOperationsByUserId - RemoveOperationByUserId: %v", err)
		}
	}
	
	return nil
}

func (s *OperationService)GetAllSlugsByUserId(ctx context.Context, userId int) ([]string, error) {
	slugs, err := s.operationRepo.GetAllSlugsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("OperationService.GetAllSlugsByUserId - GetAllSlugsByUserId: %v", err)
	}
	output := make([]string, 0, len(slugs))
	for _, id := range slugs {
		var slug entity.Slug
		slug, err := s.slugRepo.GetSlugById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("OperationService.GetAllSlugsByUserId - GetSlugById: %v", err)
		}
		output = append(output, slug.Name)
	}
	return output, nil
}
