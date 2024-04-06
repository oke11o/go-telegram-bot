package service

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

type repo interface {
	SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error)
	SaveUser(ctx context.Context, income model.User) (model.User, error)
}

func NewIncomeServce(repo repo) *IncomeService {
	return &IncomeService{
		repo: repo,
	}
}

type IncomeService struct {
	repo repo
}

func (s *IncomeService) Income(ctx context.Context, requestID string, update tgbotapi.Update) (model.User, error) {
	incomeRequest, err := model.NewIncomeRequestFromTgUpdate(requestID, update)
	if err != nil {
		return model.User{}, fmt.Errorf("model.NewIncomeRequestFromTgUpdate(%s) err: %w", requestID, err)
	}
	incomeRequest, err = s.repo.SaveIncome(ctx, incomeRequest)
	if err != nil {
		return model.User{}, fmt.Errorf("repo.SaveIncome() err: %w", err)
	}

	user, err := model.NewUserFromTgUpdate(update)
	if err != nil {
		return model.User{}, fmt.Errorf("model.NewUserFromTgUpdate() err: %w", err)
	}

	user, err = s.repo.SaveUser(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("repo.SaveUser() err: %w", err)
	}

	return user, nil
}
