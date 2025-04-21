package documentstore

import (
	"errors"
	"fmt"
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
		return fmt.Errorf("missing '%s' field", c.config.PrimaryKey)
	}
	if keyField.Type != DocumentFieldTypeString {
		return ErrInvalidKeyType
	}

	key, ok := keyField.Value.(string)
	if !ok || key == "" {
		return ErrEmptyKey
	}

	c.documents[key] = doc
	return nil
}

// Get намагається отримати документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Get(key string) (*Document, error) {
	doc, ok := c.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

// Delete намагається видалити документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Delete(key string) error {
	_, ok := c.documents[key]
	if !ok {
		return ErrDocumentNotFound
	}
	delete(c.documents, key)
	return nil
}

// List повертає всі документи в колекції
func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		result = append(result, doc)
	}
	return result
}

func (s *Store) NumCollections() int {
	return len(s.collections)
}
