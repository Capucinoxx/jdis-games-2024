package manager

import (
  "errors"

  "github.com/capucinoxx/forlorn/pkg/connector"
  "github.com/capucinoxx/forlorn/pkg/utils"
  "github.com/google/uuid"
  "go.mongodb.org/mongo-driver/bson"
)


type TokenInfo struct {
  Token     string  `bson:"token"`
  Username  string  `bson:"username"`
  Color     int     `bson:"color"`
}


type UserInfo struct {
  Username  string  `bson:"username"`
  Color     int     `bson:"color"`
}


type AuthManager struct {
  service    *connector.MongoService
  collection string
}


func NewAuthManager(db *connector.MongoService) *AuthManager {
  return &AuthManager{
    service:    db,
    collection: "users",
  }
}


func (am *AuthManager) Register(username string) (string, error) {
  filter := bson.M{"username": username}

  if v, _ := am.service.FindOne(am.collection, filter); v != nil {
    return "", errors.New("user already exist")
  }

  token := am.uuid()
  user := bson.M{"username": username, "token": token, "color": utils.NameColor(username)}

  _, err := am.service.Insert(am.collection, user)
  if err != nil {
    return "", errors.New("error inserting user")
  }

  return token, nil
}


func (am *AuthManager) Authenticate(token string) (string, int, bool) {
  filter := bson.M{"token": token}
  v, err := am.service.FindOne(am.collection, filter)
  if v == nil || err != nil {
    return "", 0, false
  }

  var result TokenInfo
  if err = v.Decode(&result); err != nil {
    return "", 0, false
  }

  return result.Username, result.Color, true
}


func (am *AuthManager) uuid() string {
  return uuid.NewString()
}


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
