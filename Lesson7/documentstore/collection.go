package documentstore

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidKeyType   = errors.New("document key must be of type string")
	ErrEmptyKey         = errors.New("document key cannot be empty")
	ErrInvalidFieldType = errors.New("invalid field type")
)

type Collection struct {
	config    *CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

// Put додає документ у колекцію або повертає помилку, якщо щось пішло не так
func (c *Collection) Put(doc Document) error {
	keyField, ok := doc.Fields[c.config.PrimaryKey]
	if !ok {
		slog.Debug("Put: missing primary key field", slog.String("primaryKey", c.config.PrimaryKey), slog.Any("document", doc))
		return fmt.Errorf("missing '%s' field", c.config.PrimaryKey)
	}
	if keyField.Type != DocumentFieldTypeString {
		slog.Debug("Put: invalid key type", slog.String("primaryKey", c.config.PrimaryKey), slog.Any("keyField", keyField))
		return ErrInvalidKeyType
	}

	key, ok := keyField.Value.(string)
	if !ok || key == "" {
		slog.Debug("Put: empty key", slog.String("primaryKey", c.config.PrimaryKey), slog.Any("keyField", keyField))
		return ErrEmptyKey
	}

	if _, exists := c.documents[key]; exists {
		slog.Debug("Put: replacing existing document", slog.String("primaryKey", c.config.PrimaryKey), slog.String("key", key), slog.Any("document", doc))
	} else {
		slog.Debug("Put: adding new document", slog.String("primaryKey", c.config.PrimaryKey), slog.String("key", key), slog.Any("document", doc))
	}

	c.documents[key] = doc
	return nil
}

// Get намагається отримати документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Get(key string) (*Document, error) {
	doc, ok := c.documents[key]
	if ok {
		slog.Debug("Get: document found", slog.String("key", key), slog.Any("document", doc))
		return &doc, nil
	}
	slog.Debug("Get: document not found", slog.String("key", key))
	return nil, ErrDocumentNotFound
}

// Delete намагається видалити документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Delete(key string) error {
	if _, ok := c.documents[key]; !ok {
		slog.Debug("Delete: document not found", slog.String("key", key))
		return ErrDocumentNotFound
	}
	slog.Debug("Delete: document deleted", slog.String("key", key))
	delete(c.documents, key)
	return nil
}

// List повертає всі документи в колекції
func (c *Collection) List() []Document {
	slog.Debug("List: retrieving all documents")
	result := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		result = append(result, doc)
	}
	slog.Debug("List: retrieved documents", slog.Int("count", len(result)))
	return result
}

// NumDocuments повертає кількість документів у колекції
func (c *Collection) NumDocuments() int {
	slog.Debug("NumDocuments: getting document count", slog.Int("count", len(c.documents)))
	return len(c.documents)
}
