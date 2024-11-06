package handler

import (
	"context"
	"fmt"
	"twt/internal/dto"
	"twt/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandler struct {
	bot *tgbotapi.BotAPI
	uc  *usecase.UserUsecase
}

func NewUserHandler(bot *tgbotapi.BotAPI, uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *UserHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.Handle"
	userID := update.Message.From.ID

	status, err := h.uc.GetUserStatus(ctx, userID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	switch status {
	case dto.UserStatusName:
		h.Name(ctx, update)
	case dto.UserStatusSurname:
		h.Surname(ctx, update)
	}
}

func (h *UserHandler) Name(ctx context.Context, update tgbotapi.Update) {
	op := "UserHandler.Name"
	userID := update.Message.From.ID
	updateText := update.Message.Text

	if err := h.uc.UpdateName(ctx, userID, updateText); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *UserHandler) Surname(ctx context.Context, update tgbotapi.Update) {
	op := "UserHandler.Surname"
	userID := update.Message.From.ID
	updateText := update.Message.Text

	if err := h.uc.UpdateSurname(ctx, userID, updateText); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
