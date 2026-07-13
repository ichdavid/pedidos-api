package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Cliente struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
}

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:SUA_SENHA@localhost:5433/loja?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	http.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var c Cliente
			json.NewDecoder(r.Body).Decode(&c)

			_, err := conn.Exec(context.Background(),
				"INSERT INTO clientes (name, email, password_hash) VALUES ($1, $2, $3)",
				c.Name, c.Email, c.PasswordHash,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	})

	http.ListenAndServe(":8080", nil)
}
