package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories/userrepositorymock" // Importa el mock
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRouter(rdb *redis.Client, userRepo *userrepositorymock.MockUserRepository) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, rdb, userRepo) // Pasa el mock del repositorio
	return router
}

type CreateTweetRequest struct {
	Content string `json:"content"`
	UserID  int64  `json:"authorId"`
}

type CreateTweetResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func TestTweetCreation(t *testing.T) {
	// Mock Redis
	rdb := getMockRedis() // Puedes dejar el mock de Redis como está

	// Crear el mock de repositorio de usuarios
	userRepo := &userrepositorymock.MockUserRepository{}

	// Setup router con el repositorio mockeado
	router := setupRouter(rdb, userRepo)

	// Primero, crea un usuario con el formato correcto
	userPayload := map[string]interface{}{"name": "Mauricio Giaconia", "email": "maurigiaconia@hotmail.com", "password": "1223"}
	userRequestBody, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest("POST", "/users/create", bytes.NewBuffer(userRequestBody))

	req.Header.Set("Content-Type", "application/json")
	_ = httptest.NewRecorder()

	// Aquí ejecutas la solicitud de creación de usuario para que el usuario esté en la base de datos, aunque en el mock

	// Simulando creación de tweets
	requests := []struct {
		payload  CreateTweetRequest
		expected int    // Código de estado esperado
		message  string // Mensaje esperado en la respuesta
	}{
		{CreateTweetRequest{Content: "test posteado", UserID: 1}, int(http.StatusCreated), "Tweet posted"},
		{CreateTweetRequest{Content: "Uala ha transformado la experiencia financiera de millones de usuarios al ofrecer una plataforma accesible y completa que les permite gestionar su dinero, realizar pagos, ahorrar, solicitar créditos y acceder a una amplia gama de servicios financieros con solo unos clics, facilitando su día a día.", UserID: 1}, int(http.StatusBadRequest), "The content of the tweet must not exceed 280 characters"},
		{CreateTweetRequest{Content: "test", UserID: 9999}, int(http.StatusBadRequest), "Nonexistent user"},
	}

	for i, tc := range requests {
		t.Run(fmt.Sprintf("Request %d", i+1), func(t *testing.T) {
			// Crea el cuerpo de la solicitud
			requestBody, _ := json.Marshal(tc.payload)

			// Crea la solicitud HTTP
			req, _ := http.NewRequest("POST", "/tweets/create", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Ejecuta la solicitud
			router.ServeHTTP(w, req)

			// Valida el código de estado
			assert.Equal(t, tc.expected, w.Code)

			// Valida el mensaje de la respuesta si es aplicable
			var response CreateTweetResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.message, response.Data)
			} else {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, tc.expected, errorResponse.Code)
				assert.Contains(t, errorResponse.Error, tc.message)
			}
		})
	}
}

func getMockRedis() *redis.Client {
	// Puedes simular el comportamiento de Redis si es necesario
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Dirección de Redis para pruebas
	})
	return client
}
