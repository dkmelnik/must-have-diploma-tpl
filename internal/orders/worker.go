package orders

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/dkmelnik/go-musthave-diploma/internal/logger"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders/dto"
)

type worker struct {
	interval        time.Duration
	attempts        int
	accrualAddr     string
	orderRepository orderRepository
}

func newWorker(accrualAddr string, or orderRepository) *worker {
	return &worker{
		interval:        300 * time.Millisecond,
		attempts:        15,
		accrualAddr:     accrualAddr,
		orderRepository: or,
	}
}

func (s *worker) CalculateAccrual(number string) {
	go func() {
		defer logger.Log.Warn("calculateAccrualForOrder", "worker", "stop")

		ctx := context.Background()

		header := map[string]string{
			"Content-Type": "text/plain",
		}

		client := resty.New()
		client.R().
			SetHeaders(header)

		sendPeriod := time.NewTicker(s.interval)
		defer sendPeriod.Stop()

		currentAttempts := 0

		for {
			select {
			case <-sendPeriod.C:
				if currentAttempts == s.attempts {
					sendPeriod.Stop()
					ctx.Done()
					return
				}
				if err := s.processAccrualRequest(ctx, client, number); err != nil {
					logger.Log.Info("calculateAccrualForOrder", "error", err)
				}
				currentAttempts++
			case <-ctx.Done():
				logger.Log.Info("calculateAccrualForOrder", "context", "cancelled")
				return
			}
		}
	}()
}

func (s *worker) processAccrualRequest(ctx context.Context, client *resty.Client, number string) error {
	accrualRes := dto.Accrual{}
	resp, err := client.R().
		SetContext(ctx).
		SetResult(&accrualRes).
		Get(fmt.Sprintf("%s/api/orders/%s", s.accrualAddr, number))
	logger.Log.Info("-------------processAccrualRequest", "number", number)
	logger.Log.Info("-------------processAccrualRequest", "accrualRes", accrualRes)

	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if accrualRes.Status == string(models.OrderRegistered) {
			return nil
		}

		order := &models.Order{
			Number: number,
		}
		if accrualRes.Status == string(models.OrderProcessed) {
			order.SetAccrual(accrualRes.Accrual)
		}

		order.Status = models.OrderStatus(accrualRes.Status)

		if err := s.orderRepository.UpdateByNumber(ctx, order); err != nil {
			logger.Log.Info("processAccrualRequest", "order", order)
			return err
		}
	case http.StatusNoContent:
		logger.Log.Warn("calculateAccrualForOrder", "http.Status", http.StatusNoContent, "resp", resp.String())
		ctx.Done()
		return nil
	case http.StatusInternalServerError:
		logger.Log.Warn("calculateAccrualForOrder", "http.Status", http.StatusInternalServerError, "resp", resp.String())
		ctx.Done()
		return nil
	}
	return nil
}
