package repositories

import (
	"database/sql"
	"time"

	"github.com/GD-Solutions/gds_base_api/internal/models"
	"github.com/GD-Solutions/gds_base_api/pkg/db"
)

// Repositorio de autenticación
type AuthRepository struct{}

// Nueva instancia del repositorio de autenticación
func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

// Verifica si el token unico existe en la tabla contract
func (r *AuthRepository) IsTokenValid(authToken string) bool {
	var count int
	query := "SELECT COUNT(*) FROM contract WHERE unique_token = ?" // Se parametriza la consulta utilizando "?"
	err := db.DB.QueryRow(query, authToken).Scan(&count)            // Le asigno el valor obtenido a la variable count, si hay algun error se guarda en err

	if err != nil {
		return false
	}
	return count == 1
}

// Obtiene un contrato basado en el authToken
func (r *AuthRepository) GetContractByToken(authToken string) (*models.Contract, error) {
	var contract models.Contract
	var dateRaw string
	query := "SELECT * FROM contract WHERE unique_token = ?"

	// Se ejecuta la query y se guardan los resultados en la variable contract
	err := db.DB.QueryRow(query, authToken).Scan(
		&contract.ContractID,
		&contract.UniqueToken,
		&dateRaw,
		&contract.Status,
		&contract.OwnerDNI,
		&contract.PlanId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el contrato
		}
		return nil, err // Error en la consulta
	}

	// Convertir el valor crudo a time.Time
	if len(dateRaw) > 0 {
		parsedDate, parseErr := time.Parse("2006-01-02 15:04:05", dateRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		contract.Date = &parsedDate
	}

	return &contract, nil
}
