package handler

import (
	"context"
	"fmt"
	"twt/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	bot *tgbotapi.BotAPI
	uc  *usecase.UserUsecase
}

func NewCommandHandler(bot *tgbotapi.BotAPI, uc *usecase.UserUsecase) *CommandHandler {
	return &CommandHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		h.Start(ctx, update)
	}
}

func (h *CommandHandler) Start(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Start"

	userID := update.Message.From.ID

	err := h.uc.CreateUser(ctx, userID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
