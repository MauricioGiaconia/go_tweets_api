package repositories

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

func FollowUser(db *sql.DB, userFollow *models.UserFollow) (bool, error) {
	query := `INSERT INTO follows (follower_id, followed_id) VALUES (?, ?)`

	_, err := db.Exec(query, userFollow.FollowerID, userFollow.FollowedID)
	if err != nil {
		fmt.Println("[x] Error to create follow: %v", err)
		return false, fmt.Errorf("[x] Error to create follow: %v", err)
	}

	return true, nil
}

func GetFollowers(db *sql.DB, userId int64) (*models.UserFollowers, error) {
	query := `SELECT *
				FROM users u
				JOIN follows f ON u.id = f.followed_id
				WHERE f.follower_id = ?`

	var followers []models.User
	rows, err := db.Query(query, userId)

	if err != nil {
		return nil, fmt.Errorf("Error fetching followers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {

		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		fmt.Println(user)
		fmt.Println(err)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		followers = append(followers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", err)
	}

	userFollowers := &models.UserFollowers{
		UserID:    int(userId),
		Followers: followers,
	}

	return userFollowers, nil
}
