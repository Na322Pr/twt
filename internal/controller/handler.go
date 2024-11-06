package controller

import (
	"context"
	"fmt"
	"twt/internal/controller/handler"
	"twt/internal/dto"
	"twt/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
	bot            *tgbotapi.BotAPI
	uc             *usecase.UserUsecase
	commandHandler *handler.CommandHandler
	userHandler    *handler.UserHandler
}

func NewController(bot *tgbotapi.BotAPI, uc *usecase.UserUsecase) *Controller {
	controller := &Controller{
		bot:            bot,
		uc:             uc,
		commandHandler: handler.NewCommandHandler(bot, uc),
		userHandler:    handler.NewUserHandler(bot, uc),
	}

	return controller
}

func (c *Controller) HandleUpdates(ctx context.Context) {
	op := "Controller.HandleUpdates"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			fmt.Printf("%s: received non-message update: %+v\n", op, update)
			continue
		}

		if update.Message.IsCommand() {
			c.commandHandler.Handle(ctx, update)
			continue
		}

		status, err := c.uc.GetUserStatus(ctx, update.Message.From.ID)
		if err != nil {
			fmt.Printf("%s: %v", op, err)
		}

		switch status {
		case dto.UserStatusName, dto.UserStatusSurname:
			c.userHandler.Handle(ctx, update)
		}
	}
}
