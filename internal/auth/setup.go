package auth

import "database/sql"

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}
