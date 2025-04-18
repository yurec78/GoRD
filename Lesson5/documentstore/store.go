package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Визначення помилок
var (
	ErrCollectionAlreadyExists  = errors.New("collection already exists")
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrInvalidDataField         = errors.New("document does not contain a valid 'data' field")
	ErrInvalidDataFieldValue    = errors.New("invalid 'data' field value")
	ErrMarshalFailed            = errors.New("failed to marshal input")
	ErrUnmarshalFailed          = errors.New("failed to unmarshal into object")
	ErrUnmarshalToMapFailed     = errors.New("failed to unmarshal JSON to map")
	ErrUnsupportedDocumentField = errors.New("unsupported document field type")
)

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

// CreateCollection створює нову колекцію, якщо вона не існує
func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, exists := s.collections[name]; exists {
		return nil, ErrCollectionAlreadyExists
	}

	collection := &Collection{
		config:    cfg,
		documents: make(map[string]Document),
	}
	s.collections[name] = collection
	return collection, nil
}

// GetCollection повертає колекцію за ім'ям
func (s *Store) GetCollection(name string) (*Collection, error) {
	col, ok := s.collections[name]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return col, nil
}

// DeleteCollection видаляє колекцію за ім'ям
func (s *Store) DeleteCollection(name string) error {
	_, ok := s.collections[name]
	if !ok {
		return ErrCollectionNotFound
	}
	delete(s.collections, name)
	return nil
}

// GetAll повертає всі документи колекції
func (c *Collection) GetAll() ([]Document, error) {
	if len(c.documents) == 0 {
		return nil, ErrDocumentNotFound
	}

	docs := make([]Document, 0, len(c.documents)) // попередньо виділяємо місце
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

func MarshalDocument(input any) (*Document, error) {
	// Спочатку маршалимо вхідний об'єкт у JSON
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}

	// Розпарсимо JSON у map[string]any
	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalToMapFailed, err)
	}

	// Створюємо документ
	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	// Перетворюємо кожне поле у відповідний тип DocumentField
	for key, value := range raw {
		field, err := toDocumentField(value)
		if err != nil {
			return nil, fmt.Errorf("%w: field '%s'", ErrUnsupportedDocumentField, key)
		}
		doc.Fields[key] = field
	}

	return doc, nil
}

func toDocumentField(value any) (DocumentField, error) {
	switch v := value.(type) {
	case string:
		return DocumentField{Type: DocumentFieldTypeString, Value: v}, nil
	case float64: // JSON числа парсяться як float64
		return DocumentField{Type: DocumentFieldTypeNumber, Value: v}, nil
	case bool:
		return DocumentField{Type: DocumentFieldTypeBool, Value: v}, nil
	case []any:
		return DocumentField{Type: DocumentFieldTypeArray, Value: v}, nil
	case map[string]any:
		return DocumentField{Type: DocumentFieldTypeObject, Value: v}, nil
	default:
		return DocumentField{}, ErrUnsupportedDocumentField
	}
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
