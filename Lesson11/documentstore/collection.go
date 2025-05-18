package documentstore

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

// Передбачається, що типи DocumentFieldType, DocumentField, Document
// визначені в іншому файлі цього ж пакету documentstore.

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidKeyType   = errors.New("document key must be of type string")
	ErrEmptyKey         = errors.New("document key cannot be empty")
	ErrInvalidFieldType = errors.New("invalid field type")
	ErrIndexExists      = errors.New("index already exists")
	ErrIndexNotFound    = errors.New("index does not exist")
)

type Collection struct {
	docsMu    sync.RWMutex // М'ютекс для захисту documents
	indexesMu sync.RWMutex // М'ютекс для захисту indexes
	config    *CollectionConfig
	documents map[string]Document
	indexes   map[string]*Index
}

type CollectionConfig struct {
	PrimaryKey string
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

type Index struct {
	FieldName string
	Sorted    []indexedEntry
}

type indexedEntry struct {
	Key      string
	Document Document
}

func (c *Collection) CreateIndex(fieldName string) error {
	c.indexesMu.Lock()
	if c.indexes == nil {
		c.indexes = make(map[string]*Index)
	}
	if _, exists := c.indexes[fieldName]; exists {
		c.indexesMu.Unlock()
		return ErrIndexExists
	}

	var entriesToProcess []struct {
		Key   string
		Value string
		Doc   Document
	}

	c.docsMu.RLock()
	if c.documents != nil {
		entriesToProcess = make([]struct {
			Key   string
			Value string
			Doc   Document
		}, 0, len(c.documents))
		for k, doc := range c.documents {
			field, ok := doc.Fields[fieldName]
			if !ok || field.Type != DocumentFieldTypeString {
				continue
			}
			valStr, ok := field.Value.(string)
			if !ok {
				slog.Warn("CreateIndex: field type is string but value is not", slog.String("key", k), slog.String("fieldName", fieldName))
				continue
			}
			entriesToProcess = append(entriesToProcess, struct {
				Key   string
				Value string
				Doc   Document
			}{k, valStr, doc})
		}
	}
	c.docsMu.RUnlock()

	sort.Slice(entriesToProcess, func(i, j int) bool {
		return entriesToProcess[i].Value < entriesToProcess[j].Value
	})

	sorted := make([]indexedEntry, len(entriesToProcess))
	for i, entry := range entriesToProcess {
		sorted[i] = indexedEntry{Key: entry.Key, Document: entry.Doc}
	}

	c.indexes[fieldName] = &Index{FieldName: fieldName, Sorted: sorted}
	c.indexesMu.Unlock()
	slog.Info("CreateIndex: successfully created index", slog.String("fieldName", fieldName))
	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	c.indexesMu.Lock()
	defer c.indexesMu.Unlock()
	if c.indexes == nil {
		return ErrIndexNotFound
	}
	if _, exists := c.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}
	delete(c.indexes, fieldName)
	slog.Info("DeleteIndex: successfully deleted index", slog.String("fieldName", fieldName))
	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.indexesMu.RLock()
	index, exists := c.indexes[fieldName]
	if !exists {
		c.indexesMu.RUnlock()
		return nil, ErrIndexNotFound
	}

	indexSortedCopy := make([]indexedEntry, len(index.Sorted))
	copy(indexSortedCopy, index.Sorted)
	c.indexesMu.RUnlock()

	var result []Document
	for _, entry := range indexSortedCopy {
		valField, ok := entry.Document.Fields[fieldName]
		if !ok || valField.Type != DocumentFieldTypeString {
			continue
		}
		val, ok := valField.Value.(string)
		if !ok {
			continue
		}

		if params.MinValue != nil && val < *params.MinValue {
			continue
		}
		if params.MaxValue != nil && val > *params.MaxValue {
			continue
		}
		result = append(result, entry.Document)
	}

	if params.Desc {
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
	}
	return result, nil
}

func (c *Collection) Put(doc Document) error {
	if c.config == nil {
		return errors.New("collection config is not initialized")
	}
	keyField, ok := doc.Fields[c.config.PrimaryKey]
	if !ok {
		return fmt.Errorf("missing '%s' field", c.config.PrimaryKey)
	}
	if keyField.Type != DocumentFieldTypeString {
		return ErrInvalidKeyType
	}

	key, ok := keyField.Value.(string)
	if !ok || key == "" {
		if !ok {
			return fmt.Errorf("primary key '%s' has type string but value is not a string", c.config.PrimaryKey)
		}
		return ErrEmptyKey
	}

	c.docsMu.Lock()
	if c.documents == nil {
		c.documents = make(map[string]Document)
	}

	_, replacing := c.documents[key]
	c.documents[key] = doc
	c.docsMu.Unlock()

	if replacing {
		slog.Debug("Put: replacing existing document", slog.String("key", key))
	} else {
		slog.Debug("Put: adding new document", slog.String("key", key))
	}

	c.indexesMu.Lock()
	if c.indexes != nil {
		c.updateIndexesUnsafe(key, doc)
	}
	c.indexesMu.Unlock()

	return nil
}

func (c *Collection) Delete(key string) error {
	c.docsMu.Lock()
	if c.documents == nil {
		c.docsMu.Unlock()
		return ErrDocumentNotFound
	}
	doc, ok := c.documents[key]
	if !ok {
		c.docsMu.Unlock()
		return ErrDocumentNotFound
	}
	delete(c.documents, key)
	c.docsMu.Unlock()

	slog.Debug("Delete: removing document", slog.String("key", key))

	c.indexesMu.Lock()
	if c.indexes != nil {
		c.removeFromIndexesUnsafe(key, doc)
	}
	c.indexesMu.Unlock()

	return nil
}

func (c *Collection) updateIndexesUnsafe(key string, doc Document) {
	for _, index := range c.indexes { // ЗМІНЕНО: fieldName на _
		// Передаємо index.FieldName в updateIndexUnsafe, оскільки саме за цим полем побудований індекс.
		c.updateIndexUnsafe(index, key, doc, index.FieldName)
	}
}

func (c *Collection) updateIndexUnsafe(index *Index, key string, doc Document, fieldName string) {
	newSorted := make([]indexedEntry, 0, len(index.Sorted))
	for _, entry := range index.Sorted {
		if entry.Key != key {
			newSorted = append(newSorted, entry)
		}
	}
	index.Sorted = newSorted

	fieldData, fieldExists := doc.Fields[fieldName]
	if !fieldExists || fieldData.Type != DocumentFieldTypeString {
		slog.Debug("updateIndexUnsafe: field not suitable for string index or removed, not adding/updating entry",
			slog.String("key", key), slog.String("fieldName", fieldName))
		return
	}

	valueToInsert, ok := fieldData.Value.(string)
	if !ok {
		slog.Warn("updateIndexUnsafe: field has type string but value is not a string, not adding/updating entry",
			slog.String("key", key), slog.String("fieldName", fieldName))
		return
	}

	insertAtIndex := sort.Search(len(index.Sorted), func(i int) bool {
		entryField, entryFieldOk := index.Sorted[i].Document.Fields[fieldName]
		if !entryFieldOk || entryField.Type != DocumentFieldTypeString {
			slog.Error("updateIndexUnsafe: inconsistent data in index during sort search",
				slog.String("indexedKey", index.Sorted[i].Key), slog.String("fieldName", fieldName))
			return true
		}
		entryValStr, valOk := entryField.Value.(string)
		if !valOk {
			slog.Error("updateIndexUnsafe: inconsistent data (value not string) in index during sort search",
				slog.String("indexedKey", index.Sorted[i].Key), slog.String("fieldName", fieldName))
			return true
		}
		return entryValStr >= valueToInsert
	})

	newEntry := indexedEntry{Key: key, Document: doc}
	index.Sorted = append(index.Sorted[:insertAtIndex], append([]indexedEntry{newEntry}, index.Sorted[insertAtIndex:]...)...)
}

func (c *Collection) removeFromIndexesUnsafe(key string, doc Document) {
	for _, index := range c.indexes {
		c.removeFromIndexUnsafe(index, key)
	}
}

func (c *Collection) removeFromIndexUnsafe(index *Index, key string) {
	newSorted := make([]indexedEntry, 0, len(index.Sorted)-1)
	for _, e := range index.Sorted {
		if e.Key != key {
			newSorted = append(newSorted, e)
		}
	}
	index.Sorted = newSorted
}

func (c *Collection) Get(key string) (*Document, error) {
	c.docsMu.RLock()
	defer c.docsMu.RUnlock()
	if c.documents == nil {
		return nil, ErrDocumentNotFound
	}
	doc, ok := c.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	docCopy := Document{Fields: make(map[string]DocumentField, len(doc.Fields))}
	for k, v := range doc.Fields {
		docCopy.Fields[k] = v
	}
	return &docCopy, nil
}

func (c *Collection) List() []Document {
	c.docsMu.RLock()
	defer c.docsMu.RUnlock()
	if c.documents == nil {
		return []Document{}
	}
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docCopy := Document{Fields: make(map[string]DocumentField, len(doc.Fields))}
		for k, v := range doc.Fields {
			docCopy.Fields[k] = v
		}
		docs = append(docs, docCopy)
	}
	return docs
}

func (c *Collection) NumDocuments() int {
	c.docsMu.RLock()
	defer c.docsMu.RUnlock()
	return len(c.documents)
}

func (c *Collection) GetAll() ([]Document, error) {
	c.docsMu.RLock()
	defer c.docsMu.RUnlock()
	if c.documents == nil || len(c.documents) == 0 {
		return []Document{}, nil
	}
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docCopy := Document{Fields: make(map[string]DocumentField, len(doc.Fields))}
		for k, v := range doc.Fields {
			docCopy.Fields[k] = v
		}
		docs = append(docs, docCopy)
	}
	return docs, nil
}
