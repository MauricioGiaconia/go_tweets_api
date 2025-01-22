package models

type UserFollow struct {
	ID         int
	UserID     int // Usuario que es seguido
	FollowerID int // Usuario que sigue
}
