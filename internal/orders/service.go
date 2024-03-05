package orders

import (
	"context"
	"errors"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/logger"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders/dto"
)

type (
	orderRepository interface {
		Save(ctx context.Context, m *models.Order) error
		FindOneByNumber(ctx context.Context, number string) (*models.Order, error)
		FindByUserID(ctx context.Context, userID models.ModelID) ([]*models.Order, error)
		UpdateByNumber(ctx context.Context, order *models.Order) error
	}

	workerService interface {
		CalculateAccrual(number string)
	}

	Service struct {
		workerService   workerService
		orderRepository orderRepository
	}
)

func NewService(ws workerService, or orderRepository) *Service {
	return &Service{ws, or}
}

func (s *Service) CreateIfNotExist(ctx context.Context, userID, orderNumber string) error {
	existingOrder, finderr := s.orderRepository.FindOneByNumber(ctx, orderNumber)
	if finderr != nil {
		if !errors.Is(finderr, apperrors.ErrNotFound) {
			return finderr
		}
		newOrder := &models.Order{
			UserID: models.ModelID(userID),
			Number: orderNumber,
			Status: models.OrderNew,
		}
		if err := s.orderRepository.Save(ctx, newOrder); err != nil {
			logger.Log.Warn("createIfNotExist", "Save", err)
			return err
		}
		s.workerService.CalculateAccrual(orderNumber)
		return nil
	}
	if existingOrder.UserID == models.ModelID(userID) {
		return errors.New("order already exists for the user")
	}
	return errors.New("order already exists")
}

func (s *Service) GetAllUserOrders(ctx context.Context, userID models.ModelID) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var out []dto.OrderResponse

	for _, v := range orders {
		d := dto.OrderResponse{
			Number:     v.Number,
			Status:     string(v.Status),
			UploadedAt: v.CreatedAt,
		}
		if v.Accrual.Valid {
			d.SetAccrual(v.Accrual.Float64)
		}
		out = append(out, d)
	}

	return out, nil
}
