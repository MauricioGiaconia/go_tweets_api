package models

import "time"

type User struct {
	ID        string    `json:"id"`                 // Identificador unico del usuario
	Name      string    `json:"name"`               // Nombre
	Email     string    `json:"email"`              // Email
	Password  *string   `json:"password,omitempty"` // Password - El ideal es que este hasheada, para efectos de esta prueba no se hará esa lógica - A su vez, es un campo opcional en el modelo por razones de seguridad
	CreatedAt time.Time `json:"created_at"`         // Fecha de creación del usuario
}
