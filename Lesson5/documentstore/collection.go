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
	Config    *CollectionConfig
	Documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

// Put додає документ у колекцію або повертає помилку, якщо щось пішло не так
func (c *Collection) Put(doc Document) error {
	keyField, ok := doc.Fields[c.Config.PrimaryKey]
	if !ok {
		return fmt.Errorf("missing '%s' field", c.Config.PrimaryKey)
	}
	if keyField.Type != DocumentFieldTypeString {
		return ErrInvalidKeyType
	}

	key, ok := keyField.Value.(string)
	if !ok || key == "" {
		return ErrEmptyKey
	}

	c.Documents[key] = doc
	return nil
}

// Get намагається отримати документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Get(key string) (*Document, error) {
	doc, ok := c.Documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

// Delete намагається видалити документ за ключем, повертає помилку, якщо документ не знайдений
func (c *Collection) Delete(key string) error {
	_, ok := c.Documents[key]
	if !ok {
		return ErrDocumentNotFound
	}
	delete(c.Documents, key)
	return nil
}

// List повертає всі документи в колекції
func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.Documents))
	for _, doc := range c.Documents {
		result = append(result, doc)
	}
	return result
}
