package withdrawals

import (
	"context"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	appdto "github.com/dkmelnik/go-musthave-diploma/internal/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/withdrawals/dto"
)

type (
	withdrawalRepository interface {
		Save(ctx context.Context, order *models.Withdrawal) error
		IsEntryByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
		Find(ctx context.Context, userID models.ModelID) ([]*models.Withdrawal, error)
	}
	balanceService interface {
		GetCurrentBalance(ctx context.Context, userID models.ModelID) (appdto.Balance, error)
	}
	Service struct {
		balanceService       balanceService
		withdrawalRepository withdrawalRepository
	}
)

func NewService(bs balanceService, ws withdrawalRepository) *Service {
	return &Service{bs, ws}
}

func (s *Service) WithdrawAccrual(ctx context.Context, userID models.ModelID, d dto.WithdrawalPayload) error {

	balance, err := s.balanceService.GetCurrentBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current-d.Sum < 0 {
		return apperrors.ErrInsufficientFunds
	}

	return s.withdrawalRepository.Save(ctx, &models.Withdrawal{
		UserID:      userID,
		OrderNumber: d.Order,
		Amount:      d.Sum,
	})
}

func (s *Service) GetAllWithdrawals(ctx context.Context, userID models.ModelID) ([]dto.WithdrawalResponse, error) {
	withdrawals, err := s.withdrawalRepository.Find(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(withdrawals) == 0 {
		return nil, apperrors.ErrNoInformationAnswer
	}

	out := make([]dto.WithdrawalResponse, 0, len(withdrawals))

	for _, v := range withdrawals {
		out = append(out, dto.WithdrawalResponse{
			Order:       v.OrderNumber,
			Sum:         v.Amount,
			ProcessedAT: v.CreatedAt,
		})
	}

	return out, nil
}
