package manager

import (
	"errors"

	"github.com/capucinoxx/forlorn/pkg/connector"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Auth est une interface pour l'authentification des utilisateurs.
type Auth interface {
	Register(username string) (string, error)
	Authenticate(token string) bool
}

// AuthManager maintient une liste d'utilisateurs et de jetons d'authentification.
type AuthManager struct {
	service    *connector.MongoService
	collection string
}

// NewAuthManager crée un nouveau AuthManager.
func NewAuthManager(db *connector.MongoService) *AuthManager {
	return &AuthManager{
		service:    db,
		collection: "users",
	}
}

// Register enregistre un nouvel utilisateur et retourne un jeton d'authentification.
// Si l'utilisateur existe déjà, une erreur est retournée. Si l'enregistrement est
// réussi, le jeton d'authentification est retourné.
func (am *AuthManager) Register(username string) (string, error) {
	filter := bson.M{"username": username}

	if v, _ := am.service.FindOne(am.collection, filter); v != nil {
		return "", errors.New("user already exist")
	}

	token := am.uuid()
	user := bson.M{"username": username, "token": token}

	_, err := am.service.Insert(am.collection, user)
	if err != nil {
		return "", errors.New("error inserting user")
	}

	return token, nil
}

// Authenticate retourne vrai si le jeton d'authentification existe. Sinon, retourne faux.
func (am *AuthManager) Authenticate(token string) bool {
	filter := bson.M{"token": token}
	v, _ := am.service.FindOne(am.collection, filter)
	return v != nil
}

// uuid génère un nouvel identifiant unique universel.
func (am *AuthManager) uuid() string {
	return uuid.NewString()
}

// Users retourne une liste de tous les utilisateurs enregistrés.
func (am *AuthManager) Users() ([]string, error) {
	bsonUsers, err := am.service.FindKeep(am.collection, bson.M{}, &bson.M{"username": 1})
	if err != nil {
		return []string{}, err
	}

	users := make([]string, len(bsonUsers))
	for i, user := range bsonUsers {
		users[i] = user["username"].(string)
	}

	return users, nil
}
