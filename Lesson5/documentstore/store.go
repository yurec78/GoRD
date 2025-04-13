package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Визначення помилок
var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
	ErrInvalidDataField        = errors.New("document does not contain a valid 'data' field")
	ErrInvalidDataFieldValue   = errors.New("invalid 'data' field value")
	ErrMarshalFailed           = errors.New("failed to marshal input")
	ErrUnmarshalFailed         = errors.New("failed to unmarshal into object")
)

type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

// CreateCollection створює нову колекцію, якщо вона не існує
func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, exists := s.Collections[name]; exists {
		return nil, ErrCollectionAlreadyExists
	}

	collection := &Collection{
		Config:    cfg,
		Documents: make(map[string]Document),
	}
	s.Collections[name] = collection
	return collection, nil
}

// GetCollection повертає колекцію за ім'ям
func (s *Store) GetCollection(name string) (*Collection, error) {
	col, ok := s.Collections[name]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return col, nil
}

// DeleteCollection видаляє колекцію за ім'ям
func (s *Store) DeleteCollection(name string) error {
	_, ok := s.Collections[name]
	if !ok {
		return ErrCollectionNotFound
	}
	delete(s.Collections, name)
	return nil
}

// GetAll повертає всі документи колекції
func (c *Collection) GetAll() ([]Document, error) {
	if len(c.Documents) == 0 {
		return nil, ErrDocumentNotFound
	}
	var docs []Document
	for _, doc := range c.Documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

// MarshalDocument перетворює вхідний тип в документ
func MarshalDocument(input any) (*Document, error) {
	// Спочатку маршалимо вхідний об'єкт в JSON
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}

	// Тепер створюємо документ
	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	// Зберігаємо серіалізовані дані в полі документа "data"
	doc.Fields["data"] = DocumentField{
		Type:  DocumentFieldTypeString,
		Value: string(data),
	}

	return doc, nil
}

// UnmarshalDocument перетворює документ назад в тип структури
func UnmarshalDocument(doc *Document, output any) error {
	// Перевіряємо, чи існує поле "data" в документі
	dataField, ok := doc.Fields["data"]
	if !ok || dataField.Type != DocumentFieldTypeString {
		return fmt.Errorf("%w: missing or invalid 'data' field", ErrInvalidDataField)
	}

	// Отримуємо серіалізовані дані
	data, ok := dataField.Value.(string)
	if !ok {
		return ErrInvalidDataFieldValue
	}

	// Тепер анмаршалимо JSON в переданий об'єкт
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}

	return nil
}
