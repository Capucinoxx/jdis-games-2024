package manager

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Auth est une interface pour l'authentification des utilisateurs.
type Auth interface {
	Register(username string) (string, error)
	Authenticate(token string) bool
}

// AuthManager maintient une liste d'utilisateurs et de jetons d'authentification.
type AuthManager struct {
	users map[string]string
	uuids map[string]string
	mu    sync.RWMutex
}

// NewAuthManager crée un nouveau AuthManager.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		users: make(map[string]string),
		uuids: make(map[string]string),
	}
}

// Register enregistre un nouvel utilisateur et retourne un jeton d'authentification.
// Si l'utilisateur existe déjà, une erreur est retournée. Si l'enregistrement est
// réussi, le jeton d'authentification est retourné.
func (am *AuthManager) Register(username string) (string, error) {
	if _, ok := am.users[username]; ok {
		return "", fmt.Errorf("user already exists")
	}

	am.mu.Lock()
	token := am.uuid()
	am.users[username] = token
	am.uuids[token] = username
	am.mu.Unlock()

	return token, nil
}

// Authenticate retourne vrai si le jeton d'authentification existe. Sinon, retourne faux.
func (am *AuthManager) Authenticate(token string) bool {
	am.mu.RLock()
	_, ok := am.uuids[token]
	am.mu.RUnlock()
	return ok
}

// uuid génère un nouvel identifiant unique universel.
func (am *AuthManager) uuid() string {
	return uuid.NewString()
}

// Users retourne une liste de tous les utilisateurs enregistrés.
func (am *AuthManager) Users() []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	users := make([]string, 0, len(am.users))
	for user := range am.users {
		users = append(users, user)
	}
	return users
}
