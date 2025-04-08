package main

import (
	"fmt"
	"lesson3/documentstore"
)

func main() {
	// Створення першого документа
	doc1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key":  {Type: documentstore.DocumentFieldTypeString, Value: "user:123"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Alice"},
			"age":  {Type: documentstore.DocumentFieldTypeNumber, Value: 30},
		},
	}

	// Створення другого документа
	doc2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key":     {Type: documentstore.DocumentFieldTypeString, Value: "product:456"},
			"name":    {Type: documentstore.DocumentFieldTypeString, Value: "Laptop"},
			"price":   {Type: documentstore.DocumentFieldTypeNumber, Value: 1200.50},
			"inStock": {Type: documentstore.DocumentFieldTypeBool, Value: true},
		},
	}

	// Створення документа без поля "key"
	doc3 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"title": {Type: documentstore.DocumentFieldTypeString, Value: "Invalid Document"},
		},
	}

	// Створення документа з неправильним типом поля "key"
	doc4 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key":   {Type: documentstore.DocumentFieldTypeNumber, Value: 123},
			"value": {Type: documentstore.DocumentFieldTypeString, Value: "Some value"},
		},
	}

	doc5 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"tags": {Type: documentstore.DocumentFieldTypeArray, Value: []string{"golang", "doc"}},
			"meta": {Type: documentstore.DocumentFieldTypeObject, Value: map[string]interface{}{"author": "John"}},
		},
	}

	fmt.Println("--- Put Documents ---")
	documentstore.Put(doc1)
	documentstore.Put(doc2)
	documentstore.Put(doc3) // Спробуємо додати невалідний документ
	documentstore.Put(doc4) // Спробуємо додати документ з неправильним типом ключа
	documentstore.Put(doc5) // Спробуємо додати документ з неправильним типом ключа

	fmt.Println("\n--- Get Documents ---")
	userDoc, foundUser := documentstore.Get("user:123")
	fmt.Printf("Get 'user:123': %v, Found: %v\n", userDoc, foundUser)

	productDoc, foundProduct := documentstore.Get("product:456")
	fmt.Printf("Get 'product:456': %v, Found: %v\n", productDoc, foundProduct)

	nonExistentDoc, foundNonExistent := documentstore.Get("nonexistent:key")
	fmt.Printf("Get 'nonexistent:key': %v, Found: %v\n", nonExistentDoc, foundNonExistent)

	fmt.Println("\n--- List Documents ---")
	allDocs := documentstore.List()
	fmt.Printf("List all documents: %v\n", allDocs)

	fmt.Println("\n--- Delete Document ---")
	deleted := documentstore.Delete("user:123")
	fmt.Printf("Delete 'user:123': %v\n", deleted)

	deletedNonExistent := documentstore.Delete("user:123") // Спробуємо видалити знову
	fmt.Printf("Delete 'user:123' again: %v\n", deletedNonExistent)

	fmt.Println("\n--- List Documents After Deletion ---")
	allDocsAfterDelete := documentstore.List()
	fmt.Printf("List all documents after deletion: %v\n", allDocsAfterDelete)
}
