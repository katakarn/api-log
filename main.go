package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "postgres-service"
	port     = 5432
	user     = "admin"
	password = "password"
	dbname   = "log"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// สร้าง connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(connStr)
	// เชื่อมต่อฐานข้อมูล
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query ข้อมูล
	rows, err := db.Query("SELECT id, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// สร้าง slice เพื่อเก็บข้อมูล users
	var users []User

	// Loop ผลลัพธ์และเก็บข้อมูลใน slice
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	// ตรวจสอบ error ที่เกิดในการวนลูป
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// แปลงข้อมูลเป็น JSON และส่งไปยัง client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	// ตัวอย่างการใช้ log
	log.Println("Received a log request")

	// ตัวอย่างการใช้ fmt.Fprintln เพื่อเขียนข้อความลงใน response
	fmt.Fprintln(w, "Log request received")
}

func main() {
	log.Println("Starting server...")
	// health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	// Handle GET request for "/users"
	http.HandleFunc("/users", getUsersHandler)

	// Handle GET request for "/log"
	http.HandleFunc("/log", logHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
