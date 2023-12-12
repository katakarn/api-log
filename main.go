package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grafana/loki-client-go"
	"golang.org/x/net/context"
)

var lokiURL = "http://loki:3100"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/your-endpoint", logMiddleware(apiHandler)).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// ตอบกลับ response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next(w, r)

		// บันทึก log ใน Loki
		logData := map[string]string{
			"request_method":  r.Method,
			"request_path":    r.URL.Path,
			"response_status": "200",
		}

		// Config Loki client
		cfg := loki.Config{
			URL:        lokiURL,
			BatchWait:  5 * time.Second,
			BatchSize:  100,
			Labels:     logData,
		}

		client, err := loki.New(cfg)
		if err != nil {
			// Handle error
			return
		}

		// Create context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Log to Loki
		client.PushMessage(ctx, "API Request", startTime, logData)
	}
}
