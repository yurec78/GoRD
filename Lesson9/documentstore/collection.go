package documentstore

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidKeyType   = errors.New("document key must be of type string")
	ErrEmptyKey         = errors.New("document key cannot be empty")
	ErrInvalidFieldType = errors.New("invalid field type")
	ErrIndexExists      = errors.New("index already exists")
	ErrIndexNotFound    = errors.New("index does not exist")
)

type Collection struct {
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
	if c.indexes == nil {
		c.indexes = make(map[string]*Index)
	}
	if _, exists := c.indexes[fieldName]; exists {
		return ErrIndexExists
	}

	var sorted []indexedEntry
	for k, doc := range c.documents {
		field, ok := doc.Fields[fieldName]
		if !ok || field.Type != DocumentFieldTypeString {
			continue
		}
		sorted = append(sorted, indexedEntry{Key: k, Document: doc})
	}
	sort.Slice(sorted, func(i, j int) bool {
		vi := sorted[i].Document.Fields[fieldName].Value.(string)
		vj := sorted[j].Document.Fields[fieldName].Value.(string)
		return vi < vj
	})

	c.indexes[fieldName] = &Index{FieldName: fieldName, Sorted: sorted}
	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	if c.indexes == nil {
		return ErrIndexNotFound
	}
	if _, exists := c.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}
	delete(c.indexes, fieldName)
	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	index, exists := c.indexes[fieldName]
	if !exists {
		return nil, ErrIndexNotFound
	}

	var result []Document
	for _, entry := range index.Sorted {
		val := entry.Document.Fields[fieldName].Value.(string)
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
		return ErrEmptyKey
	}

	if c.documents == nil {
		c.documents = make(map[string]Document)
	}

	if _, exists := c.documents[key]; exists {
		slog.Debug("Put: replacing existing document", slog.String("key", key))
	} else {
		slog.Debug("Put: adding new document", slog.String("key", key))
	}

	c.documents[key] = doc
	c.updateIndexes(key, doc)
	return nil
}

func (c *Collection) Delete(key string) error {
	if c.documents == nil {
		return ErrDocumentNotFound
	}
	doc, ok := c.documents[key]
	if !ok {
		return ErrDocumentNotFound
	}
	delete(c.documents, key)
	c.removeFromIndexes(key, doc)
	return nil
}

func (c *Collection) updateIndexes(key string, doc Document) {
	if c.indexes == nil {
		return
	}
	for field, index := range c.indexes {
		fieldVal, ok := doc.Fields[field]
		if !ok || fieldVal.Type != DocumentFieldTypeString {
			continue
		}

		valueToInsert, ok := fieldVal.Value.(string)
		if !ok {
			slog.Error("updateIndexes: potential panic - field value is not a string",
				slog.String("key", key),
				slog.String("field", field),
				slog.Any("value_type", fmt.Sprintf("%T", fieldVal.Value)),
			)
			continue // Skip this index if the value is not a string
		}

		// Create a new filtered slice to build the updated index
		filtered := index.Sorted[:0]
		removed := false

		// Iterate through the existing sorted index and add elements to the filtered slice
		for _, e := range index.Sorted {
			_, ok := e.Document.Fields[field].Value.(string)
			if !ok {
				slog.Error("updateIndexes: potential panic - existing field value is not a string",
					slog.String("existing_key", e.Key),
					slog.String("field", field),
					slog.Any("value_type", fmt.Sprintf("%T", e.Document.Fields[field].Value)),
				)
				continue // Skip this entry if the existing value is not a string (shouldn't happen if index creation is correct)
			}

			if !removed && e.Key == key {
				removed = true
				continue
			}
			filtered = append(filtered, e)
		}

		// Find the insertion point for the new entry
		insertIndex := sort.Search(len(filtered), func(i int) bool {
			v, ok := filtered[i].Document.Fields[field].Value.(string)
			return ok && v >= valueToInsert
		})

		// Insert the new entry at the correct position
		filtered = append(filtered[:insertIndex], append([]indexedEntry{{Key: key, Document: doc}}, filtered[insertIndex:]...)...)

		// Update the index's Sorted slice with the new filtered and inserted data
		index.Sorted = filtered
	}
}

func (c *Collection) removeFromIndexes(key string, doc Document) {
	if c.indexes == nil {
		return
	}
	for _, index := range c.indexes {
		filtered := index.Sorted[:0]
		for _, e := range index.Sorted {
			if e.Key != key {
				filtered = append(filtered, e)
			}
		}
		index.Sorted = filtered
	}
}

func (c *Collection) Get(key string) (*Document, error) {
	if c.documents == nil {
		return nil, ErrDocumentNotFound
	}
	doc, ok := c.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func (c *Collection) List() []Document {
	if c.documents == nil {
		return []Document{}
	}
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}
	return docs
}

func (c *Collection) NumDocuments() int {
	if c.documents == nil {
		return 0
	}
	return len(c.documents)
}

func (c *Collection) GetAll() ([]Document, error) {
	if c.documents == nil {
		return nil, errors.New("no documents in collection")
	}
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}
