// store.go
package main

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	dbName string
}

func NewMongoStore(client *mongo.Client, dbName string) *MongoStore {
	return &MongoStore{client: client, dbName: dbName}
}

func (s *MongoStore) CreateMongoCollection(ctx context.Context, collName string) error {
	// Перевіряємо, чи колекція вже існує
	names, err := s.client.Database(s.dbName).ListCollectionNames(ctx, bson.M{"name": collName})
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}
	if len(names) > 0 {
		// Якщо потрібно, можна повернути помилку або nil, якщо існування не є проблемою
		return fmt.Errorf("collection '%s' already exists", collName)
	}

	// Створюємо колекцію
	err = s.client.Database(s.dbName).CreateCollection(ctx, collName)
	if err != nil {
		// Обробка помилок, якщо CreateCollection не вдалося (наприклад, через race condition, якщо інша горутина створила її)
		// Для багатьох драйверів/версій MongoDB, якщо колекція вже існує, ця команда може не повернути помилку,
		// або повернути специфічну помилку "collection already exists".
		// Оскільки ми вже перевірили вище, ця помилка може вказувати на інші проблеми.
		return fmt.Errorf("failed to create collection '%s': %w", collName, err)
	}
	return nil
}

func (s *MongoStore) ListMongoCollections(ctx context.Context) ([]string, error) {
	names, err := s.client.Database(s.dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return names, nil
}

func (s *MongoStore) DeleteMongoCollection(ctx context.Context, collName string) error {
	if err := s.client.Database(s.dbName).Collection(collName).Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop collection '%s': %w", collName, err)
	}
	return nil
}

func (s *MongoStore) PutMongoDocument(ctx context.Context, collName string, document bson.M) error {
	collection := s.client.Database(s.dbName).Collection(collName)
	id, idExists := document["_id"]

	if idExists {
		var filter bson.M
		if idStr, ok := id.(string); ok {
			objID, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				filter = bson.M{"_id": objID}
				document["_id"] = objID // Оновлюємо документ для запису з правильним типом _id
			} else {
				filter = bson.M{"_id": idStr} // Використовуємо рядок як є
			}
		} else {
			// Якщо _id не рядок (наприклад, вже ObjectID, int), використовуємо як є
			filter = bson.M{"_id": id}
		}

		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(ctx, filter, document, opts)
		if err != nil {
			return fmt.Errorf("failed to replace/upsert document: %w", err)
		}
	} else {
		_, err := collection.InsertOne(ctx, document) // MongoDB згенерує _id
		if err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
	}
	return nil
}

func getIDFilter(docID string) bson.M {
	objID, err := primitive.ObjectIDFromHex(docID)
	if err == nil {
		return bson.M{"_id": objID}
	}
	return bson.M{"_id": docID} // Якщо не ObjectID hex, використовуємо як рядок
}

func (s *MongoStore) GetMongoDocument(ctx context.Context, collName string, docID string) (bson.M, error) {
	collection := s.client.Database(s.dbName).Collection(collName)
	filter := getIDFilter(docID)

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document with id '%s' not found in collection '%s'", docID, collName)
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return result, nil
}

func (s *MongoStore) ListMongoDocuments(ctx context.Context, collName string, filter bson.M, limit, skip int64) ([]bson.M, error) {
	collection := s.client.Database(s.dbName).Collection(collName)
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []bson.M
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}
	if documents == nil {
		return []bson.M{}, nil // Повертаємо порожній зріз, якщо нічого не знайдено
	}
	return documents, nil
}

func (s *MongoStore) DeleteMongoDocument(ctx context.Context, collName string, docID string) error {
	collection := s.client.Database(s.dbName).Collection(collName)
	filter := getIDFilter(docID)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("document with id '%s' not found for deletion in collection '%s'", docID, collName)
	}
	return nil
}

func (s *MongoStore) CreateMongoIndex(ctx context.Context, collName string, fieldName string, unique bool, order int) error {
	collection := s.client.Database(s.dbName).Collection(collName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: fieldName, Value: order}},
		Options: options.Index().SetUnique(unique),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		if strings.Contains(err.Error(), "IndexKeySpecsConflict") ||
			strings.Contains(err.Error(), "IndexOptionsConflict") ||
			strings.Contains(err.Error(), "already exists with different options") {
			return fmt.Errorf("index on field '%s' may already exist with different options: %w", fieldName, err)
		}
		return fmt.Errorf("failed to create index on field '%s': %w", fieldName, err)
	}
	return nil
}

func (s *MongoStore) DeleteMongoIndex(ctx context.Context, collName string, indexName string) error {
	collection := s.client.Database(s.dbName).Collection(collName)
	_, err := collection.Indexes().DropOne(ctx, indexName)
	if err != nil {
		// MongoDB помилка "index not found" має код 27
		// mongoErr, ok := err.(mongo.CommandError); ok && mongoErr.Code == 27
		if strings.Contains(strings.ToLower(err.Error()), "index not found") { // Простіша перевірка
			return fmt.Errorf("index '%s' not found for deletion: %w", indexName, err)
		}
		return fmt.Errorf("failed to delete index '%s': %w", indexName, err)
	}
	return nil
}
