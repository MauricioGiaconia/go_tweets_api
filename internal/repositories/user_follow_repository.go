package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

func FollowUser(db *sql.DB, userFollow *models.UserFollow) (bool, error) {
	query := `INSERT INTO follows (follower_id, followed_id) VALUES ($1, $2)`

	_, err := db.Exec(query, userFollow.FollowerID, userFollow.FollowedID)
	if err != nil {
		fmt.Println("[x] Error to create follow: %v", err)
		return false, fmt.Errorf("[x] Error to create follow: %v", err)
	}

	return true, nil
}

func GetFollows(db *sql.DB, userId int64, relationType string) (*models.UserFollows, error) {
	query := `SELECT u.id, u.name, u.email, u.created_at, f.created_at AS follow_date
				FROM users u
				JOIN follows f ON u.id = f.follower_id
				WHERE f.followed_id = $1`

	// Determinar la consulta según el tipo de relación
	if relationType == "followers" {
		query = `SELECT u.id, u.name, u.email, u.created_at, f.created_at AS follow_date
				FROM users u
				JOIN follows f ON u.id = f.follower_id
				WHERE f.followed_id = $1`
	} else if relationType == "following" {
		query = `SELECT u.id, u.name, u.email, u.created_at, f.created_at AS follow_date
				FROM users u
				JOIN follows f ON u.id = f.followed_id
				WHERE f.follower_id = $1`
	} else {
		return nil, fmt.Errorf("Invalid relationType: %s", relationType)
	}

	var follows []models.UserFollowInfo
	rows, err := db.Query(query, userId)

	if err != nil {
		return nil, fmt.Errorf("Error fetching follows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var followDate time.Time
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &followDate)
		fmt.Println(user)
		fmt.Println(err)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}

		followData := models.UserFollowInfo{
			FollowUserData: user,
			FollowDate:     followDate,
		}

		follows = append(follows, followData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", err)
	}

	userFollowers := &models.UserFollows{
		UserID:     int(userId),
		Follows:    follows,
		FollowType: relationType,
	}

	return userFollowers, nil
}
