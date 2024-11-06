package dto

type UserStatus string

const (
	UserStatusName    UserStatus = "name"
	UserStatusSurname UserStatus = "surname"
	UserStatusDone    UserStatus = "done"
)

type UserDTO struct {
	ID      int64      `db:"id"`
	Name    string     `db:"name"`
	Surname string     `db:"surname"`
	Seat    int        `db:"seat"`
	Status  UserStatus `db:"status"`
}
