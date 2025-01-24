package models

import "time"

type UserFollow struct {
	FollowerID int64      `json:"followerId"` // Usuario seguidor
	FollowedID int64      `json:"followedId"` // Usuario seguido
	CreatedAt  *time.Time `json:"createdAt"`  // Fecha de seguimiento
}

type UserFollowers struct {
	UserID    int    `json:"userId"`    // Identificador del usuario
	Followers []User `json:"followers"` // Array de usuarios que tiene como seguidor
}
