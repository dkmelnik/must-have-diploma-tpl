package balance

import (
	"context"

	"github.com/dkmelnik/go-musthave-diploma/internal/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
)

type (
	orderRepository interface {
		FindSumOfAccruals(ctx context.Context, userID models.ModelID) (float64, error)
	}
	withdrawalRepository interface {
		FindSumOfAmounts(ctx context.Context, userID models.ModelID) (float64, error)
	}
	Service struct {
		withdrawalRepository withdrawalRepository
		orderRepository      orderRepository
	}
)

func NewService(ws withdrawalRepository, or orderRepository) *Service {
	return &Service{ws, or}
}

func (s *Service) GetCurrentBalance(ctx context.Context, userID models.ModelID) (dto.Balance, error) {
	out := dto.Balance{}
	ac, err := s.orderRepository.FindSumOfAccruals(ctx, userID)
	if err != nil {
		return out, err
	}
	am, err := s.withdrawalRepository.FindSumOfAmounts(ctx, userID)
	if err != nil {
		return out, err
	}
	out.Current = ac - am
	out.Withdrawn = am

	return out, nil
}
