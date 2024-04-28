package connector

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoService est un service permettant de communiquer avec une base de données
// MongoDB.
type MongoService struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoService crée un nouveau service MongoDB. Lors de la
// création du service, une connexion à la base de données est établie.
// Si la connexion échoue, une erreur est retournée.
func NewMongoService(uri, dbName string) (*MongoService, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoService{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// MongoDataDocument est une interface permettant de convertir un type en un document BSON.
// Comme MongoDB stocke les données sous forme de documents BSON, chaque type de données
// doit implémenter cette interface pour être inséré dans la base de données.
type MongoDataDocument interface {
	ToBSON() bson.M
}

// Insert permet d'insérer un document dans une collection MongoDB.
func (m *MongoService) Insert(collection string, data MongoDataDocument) (*mongo.InsertOneResult, error) {
	result, err := m.db.Collection(collection).InsertOne(context.Background(), data.ToBSON())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOne permet de rechercher un document dans une collection MongoDB. Cela
// retourne un seul résultat. Si aucun résultat n'est trouvé, la méthode retourne
// une erreur.
func (m *MongoService) FindOne(collection string, filter bson.M) (*mongo.SingleResult, error) {
	result := m.db.Collection(collection).FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

// Find permet de rechercher des documents dans une collection MongoDB. Cela
// retourne un tableau de résultats.
func (m *MongoService) Find(collection string, filter bson.M) ([]bson.M, error) {
	cursor, err := m.db.Collection(collection).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var elem bson.M
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
