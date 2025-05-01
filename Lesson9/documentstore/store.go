package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
	slog.Info("STORE CREATED", slog.String("message", "Створено нове сховище"))
	return &Store{
		collections: make(map[string]*Collection),
	}
}

// CreateCollection створює нову колекцію, якщо вона не існує
func (s *Store) CreateCollection(name string, cfg *CollectionConfig) error {
	if _, exists := s.collections[name]; exists {
		slog.Warn("COLLECTION CREATE FAILED", slog.String("name", name), slog.String("message", fmt.Sprintf("Колекція '%s' вже існує", name)))
		return ErrCollectionAlreadyExists
	}

	collection := &Collection{
		config:    cfg,
		documents: make(map[string]Document),
	}
	s.collections[name] = collection
	slog.Info("COLLECTION CREATED", slog.String("name", name), slog.String("primaryKey", cfg.PrimaryKey), slog.String("message", fmt.Sprintf("Колекція '%s' створена з первинним ключем '%s'", name, cfg.PrimaryKey)))
	return nil
}

// GetCollection повертає колекцію за її назвою
func (s *Store) GetCollection(name string) (*Collection, error) {
	collection, exists := s.collections[name]
	if !exists {
		slog.Warn("COLLECTION GET FAILED", slog.String("name", name), slog.String("message", fmt.Sprintf("Колекція '%s' не знайдена", name)))
		return nil, ErrCollectionNotFound
	}
	slog.Debug("COLLECTION GET", slog.String("name", name), slog.String("message", fmt.Sprintf("Колекція '%s' отримана", name)))
	return collection, nil
}

// DeleteCollection видаляє колекцію за її назвою
func (s *Store) DeleteCollection(name string) error {
	if _, exists := s.collections[name]; !exists {
		slog.Warn("COLLECTION DELETE FAILED", slog.String("name", name), slog.String("message", fmt.Sprintf("Колекція '%s' не знайдена", name)))
		return ErrCollectionNotFound
	}
	delete(s.collections, name)
	slog.Info("COLLECTION DELETED", slog.String("name", name), slog.String("message", fmt.Sprintf("Колекція '%s' видалена", name)))
	return nil
}

func MarshalDocument(input any) (*Document, error) {
	// Спочатку маршалимо вхідний об'єкт у JSON
	data, err := json.Marshal(input)
	if err != nil {
		slog.Error("MARSHAL DOCUMENT FAILED", slog.Any("input", input), slog.Any("error", err), slog.String("message", fmt.Sprintf("Помилка маршалінгу '%v': %v", input, err)))
		return nil, fmt.Errorf("%w: %v", ErrMarshalFailed, err)
	}
	slog.Debug("MARSHAL DOCUMENT", slog.Any("input", input), slog.String("json", string(data)), slog.String("message", fmt.Sprintf("Об'єкт '%v' успішно маршалізовано в JSON: '%s'", input, string(data))))

	// Розпарсимо JSON у map[string]any
	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	if err != nil {
		slog.Error("MARSHAL DOCUMENT FAILED", slog.String("json", string(data)), slog.Any("error", err), slog.String("message", fmt.Sprintf("Помилка демаршалінгу JSON '%s' до map: %v", string(data), err)))
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalToMapFailed, err)
	}
	slog.Debug("MARSHAL DOCUMENT", slog.String("json", string(data)), slog.String("message", fmt.Sprintf("JSON '%s' успішно демаршалізовано до map", string(data))))

	// Створюємо документ
	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	// Перетворюємо кожне поле у відповідний тип DocumentField
	for key, value := range raw {
		field, err := toDocumentField(value)
		if err != nil {
			slog.Error("MARSHAL DOCUMENT FAILED", slog.String("key", key), slog.Any("value", value), slog.Any("error", err), slog.String("message", fmt.Sprintf("Непідтримуваний тип поля '%s': %v", key, value)))
			return nil, fmt.Errorf("%w: field '%s'", ErrUnsupportedDocumentField, key)
		}
		doc.Fields[key] = field
		slog.Debug("MARSHAL DOCUMENT", slog.String("key", key), slog.String("type", fmt.Sprintf("%T", value)), slog.String("message", fmt.Sprintf("Поле '%s' з типом '%T' успішно перетворено на DocumentField", key, value)))
	}
	slog.Debug("MARSHAL DOCUMENT", slog.Any("input", input), slog.String("message", fmt.Sprintf("Документ успішно створено з '%v'", input)))
	return doc, nil
}

func toDocumentField(value any) (DocumentField, error) {
	switch v := value.(type) {
	case string:
		slog.Debug("TO DOCUMENT FIELD", slog.String("type", "string"), slog.Any("value", v))
		return DocumentField{Type: DocumentFieldTypeString, Value: v}, nil
	case float64: // JSON числа парсяться як float64
		slog.Debug("TO DOCUMENT FIELD", slog.String("type", "float64"), slog.Any("value", v))
		return DocumentField{Type: DocumentFieldTypeNumber, Value: v}, nil
	case bool:
		slog.Debug("TO DOCUMENT FIELD", slog.String("type", "bool"), slog.Any("value", v))
		return DocumentField{Type: DocumentFieldTypeBool, Value: v}, nil
	case []any:
		slog.Debug("TO DOCUMENT FIELD", slog.String("type", "[]any"), slog.Any("value", v))
		return DocumentField{Type: DocumentFieldTypeArray, Value: v}, nil
	case map[string]any:
		slog.Debug("TO DOCUMENT FIELD", slog.String("type", "map[string]any"), slog.Any("value", v))
		return DocumentField{Type: DocumentFieldTypeObject, Value: v}, nil
	default:
		slog.Error("TO DOCUMENT FIELD FAILED", slog.String("type", fmt.Sprintf("%T", value)), slog.Any("value", value), slog.String("message", fmt.Sprintf("Непідтримуваний тип значення '%T': %v", value, value)))
		return DocumentField{}, ErrUnsupportedDocumentField
	}
}

// UnmarshalDocument перетворює документ назад в тип структури
func UnmarshalDocument(doc *Document, output any) error {
	// Перевіряємо, чи існує поле "data" в документі
	dataField, ok := doc.Fields["data"]
	if !ok || dataField.Type != DocumentFieldTypeString {
		slog.Warn("UNMARSHAL DOCUMENT FAILED", slog.Any("document", doc), slog.String("message", fmt.Sprintf("Відсутнє або неваліднe поле 'data' у документі '%v'", doc)))
		return fmt.Errorf("%w: missing or invalid 'data' field", ErrInvalidDataField)
	}

	// Отримуємо серіалізовані дані
	data, ok := dataField.Value.(string)
	if !ok {
		slog.Error("UNMARSHAL DOCUMENT FAILED", slog.Any("document", doc), slog.Any("dataField", dataField), slog.String("message", fmt.Sprintf("Неваліднe значення поля 'data' у документі '%v'", doc)))
		return ErrInvalidDataFieldValue
	}
	slog.Debug("UNMARSHAL DOCUMENT", slog.Any("document", doc), slog.String("data", data), slog.String("message", fmt.Sprintf("Отримано серіалізовані дані: '%s' з документа '%v'", data, doc)))

	// Тепер анмаршалимо JSON в переданий об'єкт
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		slog.Error("UNMARSHAL DOCUMENT FAILED", slog.String("data", data), slog.String("outputType", fmt.Sprintf("%T", output)), slog.Any("error", err), slog.String("message", fmt.Sprintf("Помилка демаршалінгу JSON '%s' в об'єкт '%T': %v", data, output, err)))
		return fmt.Errorf("%w: %v", ErrUnmarshalFailed, err)
	}
	slog.Debug("UNMARSHAL DOCUMENT", slog.Any("document", doc), slog.String("outputType", fmt.Sprintf("%T", output)), slog.String("message", fmt.Sprintf("Документ '%v' успішно демаршалізовано в об'єкт '%T'", doc, output)))
	return nil
}

// NumCollections повертає кількість колекцій у сховищі.
func (s *Store) NumCollections() int {
	slog.Debug("NumCollections", slog.Int("count", len(s.collections)))
	return len(s.collections)
}
