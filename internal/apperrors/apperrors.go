package apperrors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrParse               = errors.New("parse value")
	ErrTypeNotCorrect      = errors.New("type not correct")
	ErrIsExist             = errors.New("is exist")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNoRequiredValue     = errors.New("no required value")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrNoInformationAnswer = errors.New("no information to answer")
)
