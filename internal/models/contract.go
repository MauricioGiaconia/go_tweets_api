package models

import "time"

type Contract struct {
	ContractID  int        `json:"contract_id" db:"contract_id"`   // ID único del contrato
	UniqueToken string     `json:"unique_token" db:"unique_token"` // Token único para el contrato
	Date        *time.Time `json:"date" db:"date"`                 // Fecha del contrato
	Status      string     `json:"status" db:"status"`             // Estado del contrato
	OwnerDNI    string     `json:"owner_dni" db:"owner_dni"`       // DNI del propietario del contrato
	PlanId      int        `json:"plan_id" db:"plan_id`            // ID del plan contratado
}
