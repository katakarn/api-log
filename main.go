package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

var (
	elasticsearchURL = "http://elasticsearch:9200"
	indexName         = "api_logs"
)

type logEntry struct {
	RequestTime    time.Time `json:"request_time"`
	RequestMethod  string    `json:"request_method"`
	RequestPath    string    `json:"request_path"`
	ResponseStatus int       `json:"response_status"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/your-endpoint", logMiddleware(apiHandler)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// ตอบกลับ response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	// บันทึก log ใน Elasticsearch
	logData := logEntry{
		RequestTime:    time.Now(),
		RequestMethod:  r.Method,
		RequestPath:    r.URL.Path,
		ResponseStatus: http.StatusOK,
	}

	saveLogToElasticsearch(logData)
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next(w, r)

		// บันทึก log ของ request
		logData := logEntry{
			RequestTime:    startTime,
			RequestMethod:  r.Method,
			RequestPath:    r.URL.Path,
			ResponseStatus: http.StatusOK,
		}

		saveLogToElasticsearch(logData)
	}
}

func saveLogToElasticsearch(logData logEntry) {
	client, err := elastic.NewClient(elastic.SetURL(elasticsearchURL))
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()

	_, err = client.Index().
		Index(indexName).
		BodyJson(logData).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}
}
