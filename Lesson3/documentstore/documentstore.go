package documentstore

import "fmt"

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]Document{}

func Put(doc Document) {
	keyField, ok := doc.Fields["key"]
	if !ok || keyField.Type != DocumentFieldTypeString {
		fmt.Println("Error: Document must contain a 'key' field of type string.")
		return
	}

	key, ok := keyField.Value.(string)
	if !ok {
		fmt.Println("Error: 'key' field value is not a string.")
		return
	}

	documents[key] = doc
}

func Get(key string) (*Document, bool) {
	doc, ok := documents[key]
	if !ok {
		return nil, false
	}
	return &doc, true
}

func Delete(key string) bool {
	_, ok := documents[key]
	if !ok {
		return false
	}
	delete(documents, key)
	return true
}

func List() []Document {
	list := make([]Document, 0, len(documents))
	for _, doc := range documents {
		list = append(list, doc)
	}
	return list
}
