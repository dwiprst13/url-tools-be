package auth

import (
    "database/sql"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret-key-bruh")

type AuthService struct {
    DB *sql.DB
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"` 
}

func (a *AuthService) validateUser(email, password string) (int, string, bool) {
    var id int
    var hashedPassword string
    var role string

    err := a.DB.QueryRow(`SELECT id, password_hash, role FROM users WHERE email=$1`, email).
        Scan(&id, &hashedPassword, &role)
    if err != nil {
        return 0, "", false
    }

    if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
        return 0, "", false
    }

    return id, role, true
}

func generateToken(id int, email, role string) (string, error) {
    claims := &Claims{
        UserID: id,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "url-tools-be",
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func (a *AuthService) RegisterHandler(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    _, err = a.DB.Exec(`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)`,
        req.Email, string(hashedPassword), req.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (a *AuthService) LoginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    id, role, ok := a.validateUser(req.Email, req.Password)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    token, err := generateToken(id, req.Email, role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
