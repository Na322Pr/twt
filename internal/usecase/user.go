package usecase

import (
	"context"
	"fmt"
	"os"
	"twt/internal/dto"
	"twt/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserUsecase struct {
	bot     *tgbotapi.BotAPI
	repo    *repository.UserRepository
	curSeat int
	maxSeat int
}

func NewUserUsecase(bot *tgbotapi.BotAPI, repo *repository.UserRepository) *UserUsecase {
	curSeat, err := repo.GetCurrentMaxSeat(context.Background())
	if err != nil {
		curSeat = 1
	}

	return &UserUsecase{
		bot:     bot,
		repo:    repo,
		curSeat: curSeat,
		maxSeat: 70,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, userID int64) error {
	op := "UserUsecase.CreateUser"

	if err := uc.repo.CreateUser(ctx, userID, dto.UserStatusName); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msgText := "Ваше имя?"
	msg := tgbotapi.NewMessage(userID, msgText)

	if _, err := uc.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	return nil
}

func (uc *UserUsecase) UpdateName(ctx context.Context, userID int64, name string) error {
	op := "UserUsecase.UpdateName"
	if err := uc.repo.UpdateNameAndStatus(ctx, userID, name, dto.UserStatusSurname); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msgText := "Ваша фамилия?"
	msg := tgbotapi.NewMessage(userID, msgText)

	if _, err := uc.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	return nil
}

func (uc *UserUsecase) UpdateSurname(ctx context.Context, userID int64, surname string) error {
	op := "UserUsecase.UpdateSurname"
	if err := uc.repo.UpdateSurnameAndStatus(ctx, userID, surname, dto.UserStatusDone); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if uc.curSeat > uc.maxSeat {
		msgText := "Свободных мест не осталось((0("
		msg := tgbotapi.NewMessage(userID, msgText)
		if _, err := uc.bot.Send(msg); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}

	if err := uc.repo.UpdateSeat(ctx, userID, uc.curSeat); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msgText := fmt.Sprintf("Ваше место: %d", uc.curSeat)
	uc.curSeat++
	msg := tgbotapi.NewMessage(userID, msgText)

	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *UserUsecase) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "UserUsecase.UpdateSurname"
	status, err := uc.repo.GetUserStatus(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (uc *UserUsecase) GetUsersListFile(ctx context.Context, userID int64) error {
	op := "UserUsecase.GetUsersListFile"

	usersDTOs, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filePath, err := uc.writeToFile(ctx, usersDTOs)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	textFileBytes := tgbotapi.FileBytes{
		Name:  "praticipant_list.txt",
		Bytes: fileBytes,
	}

	msg := tgbotapi.NewDocument(userID, textFileBytes)
	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *UserUsecase) writeToFile(ctx context.Context, usersDTOs []dto.UserDTO) (string, error) {
	op := "UserUsecase.writeToFile"
	filePath := "static/participant_list.txt"

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}
	defer f.Close()

	_, err = f.WriteString("Имя\t\tФамилия\t\tМесто\n")
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	for _, userDTO := range usersDTOs {
		_, err := f.WriteString(
			fmt.Sprintf("%s\t%s\t%d\n",
				userDTO.Name,
				userDTO.Surname,
				userDTO.Seat),
		)
		if err != nil {
			return "", fmt.Errorf("%s: %v", op, err)
		}
	}

	return filePath, nil
}
