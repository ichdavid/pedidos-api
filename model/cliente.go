package model

import (
	"time"

	"github.com/google/uuid"
)

type Cliente struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // nunca expor na API
	CreatedAt    time.Time `json:"createdAt"`
}

// As informações do usuario enviado no corpo da requisição para criar um novo cliente
type CreateClienteRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
