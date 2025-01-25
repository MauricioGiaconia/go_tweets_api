package models

import "time"

type UserFollow struct {
	FollowerID int64      `json:"followerId"` // Usuario seguidor
	FollowedID int64      `json:"followedId"` // Usuario seguido
	CreatedAt  *time.Time `json:"createdAt"`  // Fecha de seguimiento
}

type UserFollows struct {
	UserID     int              `json:"userId"`     // Identificador del usuario
	Follows    []UserFollowInfo `json:"follows"`    // Array de usuarios que tiene como seguidores o esta siguiendo
	FollowType string           `json:"followType"` // String que indicara si el array son followers o following
}

type UserFollowInfo struct {
	FollowUserData User      `json:"followUserData"`
	FollowDate     time.Time `json:"followDate"`
}

type FollowsCache struct {
	Follows    UserFollows `json:"follows"`
	IsFullPage bool        `json:"isFullPage"`
}
