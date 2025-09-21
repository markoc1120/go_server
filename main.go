package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/markoc1120/go_server/internal/config"
	"github.com/markoc1120/go_server/internal/database"
	"github.com/markoc1120/go_server/internal/middleware"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	config         *config.Config
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>
`, cfg.fileServerHits.Load())))
}

func main() {
	const filepathRoot = "."

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	dbConn, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		config:         cfg,
	}

	appHandler := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/", middleware.MetricsInc(&apiCfg.fileServerHits)(http.StripPrefix("/app", appHandler)))

	// Health endpoint
	mux.HandleFunc("GET /api/healthz", handlerRediness)

	// User endpoints
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	// Auth endpoints
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	// Chirp endpoints
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	// Webhook endpoints
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)

	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
