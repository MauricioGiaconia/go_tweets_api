package middleware

import (
	"net/http"

	"github.com/GD-Solutions/uala_backend_challenge/internal/repositories"
	"github.com/GD-Solutions/uala_backend_challenge/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Configura el servicio de autenticación en el middleware
func AuthMiddleware(authRepo *repositories.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("auth_token")

		if authToken == "" {
			c.JSON(http.StatusUnauthorized, utils.ResponseToApi(http.StatusUnauthorized, "Auth token header is required", false, 0, 0, 0))
			c.Abort()
			return
		}

		contract, err := authRepo.GetContractByToken(authToken)

		if contract == nil || err != nil {
			c.JSON(http.StatusForbidden, utils.ResponseToApi(http.StatusForbidden, "Invalid auth token", false, 0, 0, 0))
			c.Abort()
			return
		}

		//Se guarda en la sesion los datos del contrato, esto servirá para tener un lugar de acceso rapido a los datos en caso de ser necesario
		c.Set("contractId", contract.ContractID)
		c.Set("authToken", contract.UniqueToken)
		c.Set("contractStatus", contract.Status)
		c.Set("ownerTaxId", contract.OwnerDNI)
		c.Set("planId", contract.PlanId)

		c.Next()
	}
}
