package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/m/v2/sqlc/api"
	"github.com/jackc/pgx/v5"
)

// Handlers
func testDb(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	connStr := "postgres://user:passwords@postgres:5432/API_DB?sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close(ctx)

	queries := api.New(conn)

	user, err := queries.GetUser(ctx, 1)
	if err != nil {
		http.Error(w, "An error has occured", http.StatusBadRequest)

		log.Println(err)
	}
	fmt.Println(user)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func testPublic(w http.ResponseWriter, r *http.Request) {
	log.Println("Public Test route hit!")
	io.WriteString(w, "Public route successful")
}

func testPrivate(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Protected route successful")
}

func testAuthenticated(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	url := config.HRMS.URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", config.HRMS.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	io.WriteString(w, string(body))
}
