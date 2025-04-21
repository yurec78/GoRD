package main

import (
	"Lesson6/documentstore"
	"fmt"
	"log"
)

func main() {
	// Створюємо новий Store
	store := documentstore.NewStore()

	// Використовуємо повне ім'я типу
	err := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "Data"})
	if err != nil {
		return
	}

	// Зберігаємо дамп у файл
	err = store.DumpToFile("backup.json")
	if err != nil {
		log.Fatalf("Error dumping to file: %v", err)
	}

	// Відновлюємо Store з файлу
	restored, err := documentstore.NewStoreFromFile("backup.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Відновлено %d колекцій\n", restored.NumCollections())
}
