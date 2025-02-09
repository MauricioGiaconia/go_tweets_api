package benchmark

import (
	"testing"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
)

type CreateTweetRoutineRequest struct {
	Content string `json:"content"`
	UserID  int64  `json:"authorId"`
}

func BenchmarkGetRoutineTimeline(b *testing.B) {
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
	tweetPayload := CreateTweetRoutineRequest{Content: "Tweet para obtener timeline y realizar el benchmark", UserID: 1}
	makeRequest(b, "POST", "/tweets/create", tweetPayload, router)

	// Benchmarking del endpoint de obtención de timeline
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeRequest(b, "GET", "/tweets/3/routine_timeline", nil, router)
	}
}
