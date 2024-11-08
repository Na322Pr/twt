package repository

import (
	"context"
	"database/sql"
	"fmt"
	"twt/internal/dto"
	"twt/pkg/postgres"
)

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (r *UserRepository) CreateUser(ctx context.Context, userID int64, status dto.UserStatus) error {
	op := "UserRepository.CreateUser"
	query := `INSERT INTO users(id, status) VALUES($1, $2)`

	_, err := r.Conn.Exec(ctx, query, userID, status)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateNameAndStatus(ctx context.Context, userID int64, name string, status dto.UserStatus) error {
	op := "UserRepository.UpdateNameAndStatus"
	query := `UPDATE users SET name = $2, status = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, name, dto.UserStatusSurname)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateSurnameAndStatus(ctx context.Context, userID int64, surname string, status dto.UserStatus) error {
	op := "UserRepository.UpdateSurnameAndStatus"
	query := `UPDATE users SET surname = $2, status = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, surname, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateKKAndStatus(ctx context.Context, userID int64, isKK bool, status dto.UserStatus) error {
	op := "UserRepository.UpdateSurnameAndStatus"
	query := `UPDATE users SET is_kk = $2, status = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, isKK, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateSeat(ctx context.Context, userID int64, seat int) error {
	op := "UserRepository.UpdateSeat"
	query := `UPDATE users SET seat = $2 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, seat)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "UserRepository.GetUserStatus"
	query := `SELECT status FROM users WHERE id = $1`

	var status dto.UserStatus
	err := r.Conn.QueryRow(ctx, query, userID).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (r *UserRepository) GetUserSeat(ctx context.Context, userID int64) (int, error) {
	op := "UserRepository.GetUserSeat"
	query := `SELECT seat FROM users WHERE id = $1`

	var seat int
	err := r.Conn.QueryRow(ctx, query, userID).Scan(&seat)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return seat, nil
}

func (r *UserRepository) GetCurrentMaxSeat(ctx context.Context) (int, error) {
	op := "UserRepository.GetCurrentMaxSeat"
	query := `SELECT MAX(seat) FROM users`

	var seat int
	err := r.Conn.QueryRow(ctx, query).Scan(&seat)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return seat, nil
}

func (r *UserRepository) GetUsersWithSeats(ctx context.Context) ([]dto.UserDTO, error) {
	op := "UserRepository.GetUsersWithSeats"
	query := "SELECT name, surname, seat, is_kk FROM users WHERE is_kk IS NOT NULL"

	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	usersDTOs := make([]dto.UserDTO, 0, 50)

	for rows.Next() {
		var userDTO dto.UserDTO
		if err := rows.Scan(&userDTO.Name, &userDTO.Surname, &userDTO.Seat, &userDTO.IsKK); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		usersDTOs = append(usersDTOs, userDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usersDTOs, nil
}
