package domain

import (
	"time"
)

type SigningKey struct {
	ID            string
	KID           string
	Alg           string
	PublicKeyPem  string
	PrivateKeyPem string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
