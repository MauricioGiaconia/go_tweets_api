package models

import "time"

type Tweet struct {
	ID         int64     `json:"tweetId"`    // Identificador unico del tweet
	UserID     int64     `json:"authorId"`   // Identificador del usuario creador del tweet
	AuthorName *string   `json:"authorName"` // Campo opcional: Nombre del usuario creador del tweet
	Content    string    `json:"content"`    // Contenido del tweet
	CreatedAt  time.Time `json:"createdAt"`  // Fecha de creación
}
