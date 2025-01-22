package models

import "time"

type User struct {
	ID        string    `json:"id"`         // Identificador unico del usuario
	Name      string    `json:"name"`       // Nombre
	Email     string    `json:"email"`      // Email
	Following []string  `json:"following"`  // Array de ids de los usuarios al que sigue
	CreatedAt time.Time `json:"created_at"` // Fecha de creación del usuario
}
