package manager

// AuthManager handles user authentication and registration using a MongoDB service.
// This manager provides functionality for registering new users, authenticating users
// based on a token, and retrieving a list of registered users. It ensures that user
// information is securely stored and efficiently retrieved from the database.

import (
	"errors"

	"github.com/capucinoxx/jdis-games-2024/pkg/connector"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// TokenInfo represents the structure of a user's token information stored in MongoDB.
type TokenInfo struct {
	Token    string `bson:"token" json:"token"`
	Username string `bson:"username" json:"username"`
	Color    int    `bson:"color" json:"color"`
	IsAdmin  bool   `bson:"is_admin" json:"is_admin"`
}

// UserInfo represents the structure of a user's basic information retrieved from MongoDB.
type UserInfo struct {
	Username string `bson:"username"`
	Color    int    `bson:"color"`
	IsAdmin  bool   `bson:"is_admin"`
}

// AuthManager handles user authentication and registration.
type AuthManager struct {
	service    *connector.MongoService
	collection string
}

// NewAuthManager creates a new AuthManager with the specified MongoDB service.
func NewAuthManager(db *connector.MongoService) *AuthManager {
	return &AuthManager{
		service:    db,
		collection: "users",
	}
}

// Register registers a new user with the specified username. It generates a unique token
// for the user and stores their information in the MongoDB collection. If the user already
// exists, an error is returned.
func (am *AuthManager) Register(username string, isAdmin bool) (string, error) {
	if len(username) > 16 || len(username) < 3 {
		return "", errors.New("username must be between 3 and 16 characters")
	}

	filter := bson.M{"username": username}

	if v, _ := am.service.FindOne(am.collection, filter); v != nil {
		return "", errors.New("user already exist")
	}

	token := am.uuid()
	user := bson.M{"username": username, "token": token, "color": utils.NameColor(username), "is_admin": isAdmin}

	_, err := am.service.Insert(am.collection, user)
	if err != nil {
		return "", errors.New("error inserting user")
	}

	return token, nil
}

func (am *AuthManager) List() ([]TokenInfo, error) {
	v, err := am.service.Find(am.collection, bson.M{})
	if err != nil {
		return nil, err
	}

	result := make([]TokenInfo, 0, len(v))
	for _, e := range v {
		bytes, err := bson.Marshal(e)
		if err != nil {
			continue
		}

		var res TokenInfo
		if err = bson.Unmarshal(bytes, &res); err != nil {
			continue
		}

		result = append(result, res)
	}

	return result, nil
}

// Authenticate authenticates a user based on their token. It retrieves the user's information
// from the MongoDB collection and returns their username, color, and a boolean indicating success.
func (am *AuthManager) Authenticate(token string) (string, int, bool, bool) {
	filter := bson.M{"token": token}
	v, err := am.service.FindOne(am.collection, filter)
	if v == nil || err != nil {
		return "", 0, false, false
	}

	var result TokenInfo
	if err = v.Decode(&result); err != nil {
		return "", 0, false, false
	}

	return result.Username, result.Color, result.IsAdmin, true
}

// uuid generates a new unique identifier for user tokens.
func (am *AuthManager) uuid() string {
	return uuid.NewString()
}

// Users retrieves a list of all registered users with their username and color.
func (am *AuthManager) Users() ([]UserInfo, error) {
	bsonUsers, err := am.service.FindKeep(am.collection, bson.M{}, &bson.M{"username": 1, "color": 1})
	if err != nil {
		return []UserInfo{}, err
	}

	users := make([]UserInfo, len(bsonUsers))
	for i, user := range bsonUsers {
		users[i] = UserInfo{Username: user["username"].(string), Color: user["color"].(int)}
	}

	return users, nil
}

func (am *AuthManager) SetupAdmins(admins []TokenInfo) {
	count := 0
	for _, admin := range admins {
		filter := bson.M{"username": admin.Username}

		if v, _ := am.service.FindOne(am.collection, filter); v != nil {
			continue
		}

		user := bson.M{"username": admin.Username, "token": admin.Token, "color": admin.Color, "is_admin": true}
		_, err := am.service.Insert(am.collection, user)
		if err == nil {
			count++
		}
	}

	utils.Log("config", "admins", "%d admins have been configured", count)
}
