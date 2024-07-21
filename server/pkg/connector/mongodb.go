package connector

// MongoService provides functionality for communicating with a MongoDB database.
// This service allows for connecting to the database, inserting documents, querying collections,
// and updating documents. It simplifies the interaction with MongoDB by providing easy-to-use
// methods for common database operations.
//
// This service ensures thread safety and efficient resource management by properly handling
// database connections and operations.
//
// Usage of this package involves creating an instance of MongoService with the MongoDB URI
// and database name. The MongoService instance can then be used to perform various database
// operations.

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoService provides functionality for communicating with a MongoDB database.
type MongoService struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoService creates a new MongoDB service. Upon creation, it establishes a connection
// to the database. If the connection fails, an error is returned.
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

// Insert inserts a document into a MongoDB collection.
func (m *MongoService) Insert(collection string, data bson.M) (*mongo.InsertOneResult, error) {
	result, err := m.db.Collection(collection).InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOne searches for a document in a MongoDB collection. It returns a single result.
// If no result is found, an error is returned.
func (m *MongoService) FindOne(collection string, filter bson.M) (*mongo.SingleResult, error) {
	result := m.db.Collection(collection).FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

// Find searches for documents in a MongoDB collection. It returns an array of results.
func (m *MongoService) Find(collection string, filter bson.M) ([]bson.M, error) {
	return m.FindKeep(collection, filter, nil)
}

// FindKeep searches for documents in a MongoDB collection with optional field projection.
// It returns an array of results.
func (m *MongoService) FindKeep(collection string, filter bson.M, fields *bson.M) ([]bson.M, error) {
	findOptions := options.Find()
	if fields != nil {
		findOptions.SetProjection(*fields)
	}

	cursor, err := m.db.Collection(collection).Find(context.Background(), filter, findOptions)
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

// Push adds an element to an array field in a document identified by UUID in a MongoDB collection.
func (m *MongoService) Push(collection string, uuid string, arrayField string, element interface{}) error {
	filter := bson.M{"_id": uuid}
	update := bson.M{"$push": bson.M{arrayField: element}}

	_, err := m.db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	return err
}
