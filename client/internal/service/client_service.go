package service

import (
	repo "client/internal/service/repository"
	shared "client/shared/models"
	"context"
)

type ClientService interface {
	ServiceBase[shared.Client]
}

type clientService struct {
	repo repo.ClientRepository
}

func NewClientService(repo repo.RepositoryBase[shared.Client]) ClientService {
	return &clientService{repo}
}

func (s *clientService) Save(ctx context.Context, client shared.Client) (*shared.Client, error) {
	var result *shared.Client
	var err error

	if client.Id <= 0 {
		result, err = s.repo.Create(ctx, &client)
	} else {
		result, err = s.repo.Update(ctx, &client)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *clientService) Delete(ctx context.Context, id int) (bool, error) {
	deleted, err := s.repo.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	return deleted, nil
}

func (s *clientService) Get(ctx context.Context, id int) (*shared.Client, error) {
	client, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return client, nil
}
