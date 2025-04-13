package documentstore2

import (
	"fmt"
)

type Collection struct {
	config    *CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (c *Collection) Put(doc Document) {
	keyField, ok := doc.Fields[c.config.PrimaryKey]
	if !ok || keyField.Type != DocumentFieldTypeString {
		fmt.Printf("Error: Document must contain a '%s' field of type string.\n", c.config.PrimaryKey)
		return
	}

	key, ok := keyField.Value.(string)
	if !ok || key == "" {
		fmt.Println("Error: 'key' field value is not a valid non-empty string.")
		return
	}
	c.documents[key] = doc
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, ok := c.documents[key]
	if !ok {
		return nil, false
	}
	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	_, ok := c.documents[key]
	if !ok {
		return false
	}
	delete(c.documents, key)
	return true
}

func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		result = append(result, doc)
	}
	return result
}
