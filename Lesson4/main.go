package main

import (
	"Lesson4/documentstore"
	"fmt"
)

func main() {
	// Створюємо сховище
	store := documentstore2.NewStore()

	// Створюємо колекцію "users" з primary key "key"
	created, users := store.CreateCollection("users", &documentstore2.CollectionConfig{PrimaryKey: "key"})
	if !created {
		fmt.Println("Collection 'users' already exists")
		return
	}

	// Додаємо документи в колекцію
	doc1 := documentstore2.Document{
		Fields: map[string]documentstore2.DocumentField{
			"key":  {Type: documentstore2.DocumentFieldTypeString, Value: "user:1"},
			"name": {Type: documentstore2.DocumentFieldTypeString, Value: "Alice"},
			"age":  {Type: documentstore2.DocumentFieldTypeNumber, Value: 30},
		},
	}
	doc2 := documentstore2.Document{
		Fields: map[string]documentstore2.DocumentField{
			"key":   {Type: documentstore2.DocumentFieldTypeString, Value: "user:2"},
			"name":  {Type: documentstore2.DocumentFieldTypeString, Value: "Bob"},
			"admin": {Type: documentstore2.DocumentFieldTypeBool, Value: true},
		},
	}

	users.Put(doc1)
	users.Put(doc2)

	// Отримуємо документ
	doc, found := users.Get("user:1")
	fmt.Println("--- Get user:1 ---")
	if found {
		fmt.Printf("Found: %+v\n", doc)
	} else {
		fmt.Println("Not found")
	}

	// Виводимо всі документи
	fmt.Println("\n--- List Users ---")
	all := users.List()
	for _, d := range all {
		fmt.Printf("%+v\n", d)
	}

	// Видаляємо документ
	deleted := users.Delete("user:2")
	fmt.Printf("\nDeleted user:2? %v\n", deleted)

	// Після видалення
	fmt.Println("\n--- Users After Deletion ---")
	all = users.List()
	for _, d := range all {
		fmt.Printf("%+v\n", d)
	}
}
