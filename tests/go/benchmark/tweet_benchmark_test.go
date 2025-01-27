package benchmark

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CreateTweetRequest struct {
	Content string `json:"content"`
	UserID  int64  `json:"authorId"`
}

func BenchmarkTweetCreation(b *testing.B) {
	// Configuración de la base de datos y el mock de Redis
	db, err := factory.GetDatabase("postgres")
	if err != nil {
		b.Fatalf("failed to create DB: %v", err)
	}

	rdb := getMockBenchmarkRedis()

	conn, err := db.Connect()
	if err != nil {
		b.Fatalf("failed to connect DB: %v", err)
	}
	defer conn.Close()

	router := setupTweetRouter(conn, rdb)

	// Crear un usuario de prueba
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}
	makeRequest(b, "POST", "/users/create", userPayload, router)

	// Benchmarking del endpoint de creación de tweets
	b.ResetTimer() // Restablece el temporizador para evitar medir el tiempo de configuración
	for i := 0; i < b.N; i++ {
		tweetPayload := CreateTweetRequest{Content: "Tweet para benchmark", UserID: 1}
		makeRequest(b, "POST", "/tweets/create", tweetPayload, router)
	}
}

func BenchmarkGetTimeline(b *testing.B) {
	// Configuración de la base de datos y el mock de Redis
	db, err := factory.GetDatabase("postgres")
	if err != nil {
		b.Fatalf("failed to create DB: %v", err)
	}
	rdb := getMockBenchmarkRedis()

	conn, err := db.Connect()
	if err != nil {
		b.Fatalf("failed to connect DB: %v", err)
	}
	defer conn.Close()

	router := setupTweetRouter(conn, rdb)

	// Crear un usuario de prueba
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}
	makeRequest(b, "POST", "/users/create", userPayload, router)

	// Crear un tweet
	tweetPayload := CreateTweetRequest{Content: "Tweet para obtener timeline y realizar el benchmark", UserID: 1}
	makeRequest(b, "POST", "/tweets/create", tweetPayload, router)

	// Benchmarking del endpoint de obtención de timeline
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeRequest(b, "GET", "/tweets/3/timeline", nil, router)
	}
}

func makeRequest(b *testing.B, method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			b.Fatalf("Error marshalling request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func getMockBenchmarkRedis() *redis.Client {

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

func setupTweetRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}
