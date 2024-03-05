package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	appdto "github.com/dkmelnik/go-musthave-diploma/internal/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/users/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/utils"
)

type (
	UserRepository interface {
		Save(ctx context.Context, m *models.User) (models.ModelID, error)
		IsEntryByLogin(ctx context.Context, login string) (bool, error)
		FindOneByLogin(ctx context.Context, login string) (*models.User, error)
	}
	JWTService interface {
		BuildJWTString(userID models.ModelID, jti string) (string, error)
		ParseToken(tokenString string) (*appdto.Claims, error)
	}
	Service struct {
		cost           int
		jwtService     JWTService
		userRepository UserRepository
	}
)

func NewService(jwtService JWTService, userRepository UserRepository) *Service {
	return &Service{bcrypt.DefaultCost, jwtService, userRepository}
}

func (s *Service) Register(ctx context.Context, dto dto.RegisterPayload) (string, error) {
	exist, err := s.userRepository.IsEntryByLogin(ctx, dto.Login)
	if err != nil {
		return "", err
	}

	if exist {
		return "", apperrors.ErrIsExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), s.cost)
	if err != nil {
		return "", err
	}

	userID, err := s.userRepository.Save(ctx, &models.User{
		Login:    dto.Login,
		Password: string(hashedPassword),
	})

	if err != nil {
		return "", err
	}

	token, err := s.jwtService.BuildJWTString(userID, utils.GenerateGUID())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Authenticate(ctx context.Context, dto dto.LoginPayload) (string, error) {

	user, err := s.userRepository.FindOneByLogin(ctx, dto.Login)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := s.jwtService.BuildJWTString(user.ID, utils.GenerateGUID())
	if err != nil {
		return "", err
	}

	return token, nil
}
