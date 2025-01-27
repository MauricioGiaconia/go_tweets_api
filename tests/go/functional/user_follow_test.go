package functional

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type UserFollows struct {
	UserID     int                     `json:"userId"`     // Identificador del usuario
	Follows    []models.UserFollowInfo `json:"follows"`    // Array de usuarios que tiene como seguidores o esta siguiendo
	FollowType string                  `json:"followType"` // O usa `time.Time` si quieres manejar fechas como objetos
}

type FollowsResponse struct {
	Code   int         `json:"code"`
	Data   UserFollows `json:"data"`
	Count  int         `json:"count"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Next   string      `json:"next"`
}

type FollowCreationRequest struct {
	FollowerId int64 `json:"followerId"`
	FollowedId int64 `json:"followedId"`
}

type CreateFollowResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func TestPostAndGetFollows(t *testing.T) {
	db, err := factory.GetDatabase("sqlite")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	rdb := getMockFollowRedis()

	conn, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}

	defer conn.Close()

	router := setupFollowRouter(conn, rdb)

	// Se crea un user para realizar el test de follow
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}

	w := makeFollowRequest(t, "POST", "/users/create", userPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	// Crear el segundo usuario
	otherUserPayload := map[string]interface{}{
		"name":     "Juan Perez",
		"email":    "juanperez@hotmail.com",
		"password": "4567",
	}
	w = makeFollowRequest(t, "POST", "/users/create", otherUserPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	requests := []struct {
		payload  FollowCreationRequest
		expected int    // Código de estado esperado
		message  string // Mensaje esperado en la respuesta
	}{
		{FollowCreationRequest{
			FollowerId: 1,
			FollowedId: 2,
		}, http.StatusCreated, "Followed"},
		{FollowCreationRequest{
			FollowerId: 2,
			FollowedId: 1,
		}, http.StatusCreated, "Followed"},
		{FollowCreationRequest{
			FollowerId: 1,
			FollowedId: 200,
		}, http.StatusNotFound, "Nonexistent followed ID user"},
		{FollowCreationRequest{
			FollowerId: 100,
			FollowedId: 2,
		}, http.StatusNotFound, "Nonexistent follower ID user"},
		{FollowCreationRequest{
			FollowerId: 1,
			FollowedId: 1,
		}, http.StatusBadRequest, "Cannot follow yourself"},
	}

	for i, tc := range requests {
		t.Run(fmt.Sprintf("Request %d", i+1), func(t *testing.T) {

			requestBody, _ := json.Marshal(tc.payload)

			req, _ := http.NewRequest("POST", "/users_follow/create", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, int64(tc.expected), int64(w.Code))

			var response CreateFollowResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.message, response.Data)
			} else {
				var errorResponse utils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, int64(tc.expected), int64(errorResponse.Code))
				assert.Contains(t, errorResponse.Error, tc.message)
			}
		})
	}

	w = makeFollowRequest(t, "GET", "/users_follow/1/follows/followers", nil, router)
	assert.Equal(t, int64(http.StatusOK), int64(w.Code), "Expected status OK for GET request")

	var followResponse FollowsResponse

	err = json.Unmarshal(w.Body.Bytes(), &followResponse)
	assert.NoError(t, err, "Error deserializando la respuesta GET")
	assert.True(t, len(followResponse.Data.Follows) > 0, "Expected at least one follower")

	w = makeFollowRequest(t, "GET", "/users_follow/1/follows/following", nil, router)
	assert.Equal(t, int64(http.StatusOK), int64(w.Code), "Expected status OK for GET request")

	err = json.Unmarshal(w.Body.Bytes(), &followResponse)
	assert.NoError(t, err, "Error deserializando la respuesta GET")
	assert.True(t, len(followResponse.Data.Follows) > 0, "Expected at least one following")
}

func setupFollowRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}

func getMockFollowRedis() *redis.Client {

	var ctx = context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("error conectando a Redis: %v\n", err)
		return nil
	}

	return client
}

// Funcion utilziada para realizar requests necesarias para el test (por ejemplo, si se necesita crear un usuario para poder testear los endpoints de tweets)
func makeFollowRequest(t *testing.T, method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var requestBody []byte

	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Error marshalling request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w
}
