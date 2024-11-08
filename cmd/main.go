package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"twt/internal/config"
	"twt/internal/controller"
	"twt/internal/repository"
	"twt/internal/usecase"
	"twt/pkg/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Printf(".env file not found: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.MustLoad()
	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(cfg.TG.Token)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "load", Description: "Выгрузить список участников"},
	}

	cmgCfg := tgbotapi.NewSetMyCommands(commands...)
	_, err = bot.Request(cmgCfg)
	if err != nil {
		log.Fatalf("failed to set commands: %v", err)
	}

	pg, err := postgres.Connection(cfg.PG.URL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pg.Close()

	repository := repository.NewUserRepository(pg)
	usecase := usecase.NewUserUsecase(bot, repository, cfg.SubAdminIDs)
	controller := controller.NewController(bot, usecase)

	go func() {
		controller.HandleUpdates(ctx)
	}()

	if err := NotifyOnStartUp(bot, *cfg); err != nil {
		log.Printf("failed to notify admins: %v", err)
	}

	<-stop
	log.Println("\nShutting down...")
	os.Exit(0)
}

func NotifyOnStartUp(bot *tgbotapi.BotAPI, cfg config.Config) error {
	op := "NotifyOnStartUp"

	msgText := "Бот запущен"

	for _, adminID := range cfg.TG.AdminIDs {
		msg := tgbotapi.NewMessage(adminID, msgText)

		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
