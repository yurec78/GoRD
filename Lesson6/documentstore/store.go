package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	log.Println("STORE CREATED: Створено нове сховище")
	return &Store{
		collections: make(map[string]*Collection),
	}
}

// CreateCollection створює нову колекцію, якщо вона не існує
func (s *Store) CreateCollection(name string, cfg *CollectionConfig) error {
	if _, exists := s.collections[name]; exists {
		log.Printf("COLLECTION CREATE FAILED: Колекція '%s' вже існує", name)
		return ErrCollectionAlreadyExists
	}

	collection := &Collection{
		config:    cfg,
		documents: make(map[string]Document),
	}
	s.collections[name] = collection
	log.Printf("COLLECTION CREATED: Колекція '%s' створена з первинним ключем '%s'", name, cfg.PrimaryKey)
	return nil
}

// GetCollection повертає колекцію за її назвою
func (s *Store) GetCollection(name string) (*Collection, error) {
	collection, exists := s.collections[name]
	if !exists {
		log.Printf("COLLECTION GET FAILED: Колекція '%s' не знайдена", name)
		return nil, ErrCollectionNotFound
	}
	log.Printf("COLLECTION GET: Колекція '%s' отримана", name)
	return collection, nil
}

// DeleteCollection видаляє колекцію за її назвою
func (s *Store) DeleteCollection(name string) error {
	if _, exists := s.collections[name]; !exists {
		log.Printf("COLLECTION DELETE FAILED: Колекція '%s' не знайдена", name)
		return ErrCollectionNotFound
	}
	delete(s.collections, name)
	log.Printf("COLLECTION DELETED: Колекція '%s' видалена", name)
	return nil
}

// GetAll повертає всі документи колекції
func (c *Collection) GetAll() ([]Document, error) {
	if len(c.documents) == 0 {
		log.Printf("COLLECTION GET ALL FAILED: Колекція '%s' порожня", c.config.PrimaryKey)
		return nil, ErrDocumentNotFound
	}

	docs := make([]Document, 0, len(c.documents)) // попередньо виділяємо місце
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}
	log.Printf("COLLECTION GET ALL: Отримано %d документів з колекції '%s'", len(docs), c.config.PrimaryKey)
	return docs, nil
}

func MarshalDocument(input any) (*Document, error) {
	// Спочатку маршалимо вхідний об'єкт у JSON
	data, err := json.Marshal(input)
	if err != nil {
		log.Printf("MARSHAL DOCUMENT FAILED: Помилка маршалінгу '%v': %v", input, err)
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}
	log.Printf("MARSHAL DOCUMENT: Об'єкт '%v' успішно маршалізовано в JSON: '%s'", input, string(data))

	// Розпарсимо JSON у map[string]any
	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	if err != nil {
		log.Printf("MARSHAL DOCUMENT FAILED: Помилка демаршалінгу JSON '%s' до map: %v", string(data), err)
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalToMapFailed, err)
	}
	log.Printf("MARSHAL DOCUMENT: JSON '%s' успішно демаршалізовано до map", string(data))

	// Створюємо документ
	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	// Перетворюємо кожне поле у відповідний тип DocumentField
	for key, value := range raw {
		field, err := toDocumentField(value)
		if err != nil {
			log.Printf("MARSHAL DOCUMENT FAILED: Непідтримуваний тип поля '%s': %v", key, value)
			return nil, fmt.Errorf("%w: field '%s'", ErrUnsupportedDocumentField, key)
		}
		doc.Fields[key] = field
		log.Printf("MARSHAL DOCUMENT: Поле '%s' з типом '%T' успішно перетворено на DocumentField", key, value)
	}
	log.Printf("MARSHAL DOCUMENT: Документ успішно створено з '%v'", input)
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
		log.Printf("TO DOCUMENT FIELD FAILED: Непідтримуваний тип значення '%T': %v", value, value)
		return DocumentField{}, ErrUnsupportedDocumentField
	}
}

// UnmarshalDocument перетворює документ назад в тип структури
func UnmarshalDocument(doc *Document, output any) error {
	// Перевіряємо, чи існує поле "data" в документі
	dataField, ok := doc.Fields["data"]
	if !ok || dataField.Type != DocumentFieldTypeString {
		log.Printf("UNMARSHAL DOCUMENT FAILED: Відсутнє або неваліднe поле 'data' у документі '%v'", doc)
		return fmt.Errorf("%w: missing or invalid 'data' field", ErrInvalidDataField)
	}

	// Отримуємо серіалізовані дані
	data, ok := dataField.Value.(string)
	if !ok {
		log.Printf("UNMARSHAL DOCUMENT FAILED: Неваліднe значення поля 'data' у документі '%v'", doc)
		return ErrInvalidDataFieldValue
	}
	log.Printf("UNMARSHAL DOCUMENT: Отримано серіалізовані дані: '%s' з документа '%v'", data, doc)

	// Тепер анмаршалимо JSON в переданий об'єкт
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		log.Printf("UNMARSHAL DOCUMENT FAILED: Помилка демаршалінгу JSON '%s' в об'єкт '%T': %v", data, output, err)
		return fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}
	log.Printf("UNMARSHAL DOCUMENT: Документ '%v' успішно демаршалізовано в об'єкт '%T'", doc, output)
	return nil
}
