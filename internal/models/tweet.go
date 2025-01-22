package models

import "time"

type Tweet struct {
	ID        string    `json:"id"`         // Identificador unico del tweet
	UserID    string    `json:"user_id"`    // Identificador del usuario creador del tweet
	Content   string    `json:"content"`    // Contenido del tweet
	CreatedAt time.Time `json:"created_at"` // Fecha de creación
}
