package main

import (
	"Lesson9/documentstore"
	"fmt"
	"log/slog"
	"os"
)

// User - структура для прикладу даних користувача
type User struct {
	ID    string `json:"Data"` // Використовуємо "Data" як первинний ключ відповідно до вашого CreateCollection
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Створюємо нове сховище
	store := documentstore.NewStore()

	// Створюємо колекцію користувачів з первинним ключем "Data"
	usersCollectionName := "users"
	err := store.CreateCollection(usersCollectionName, &documentstore.CollectionConfig{PrimaryKey: "Data"})
	if err != nil {
		slog.Error("Error creating users collection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Отримуємо колекцію користувачів
	usersCollection, err := store.GetCollection(usersCollectionName)
	if err != nil {
		slog.Error("Error getting users collection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Створюємо кілька користувачів
	user1 := User{ID: "user1", Name: "Alice", Email: "alice@example.com"}
	user2 := User{ID: "user2", Name: "Bob", Email: "bob@example.com"}

	// Маршалізуємо користувачів у документи та додаємо їх до колекції
	// Маршалізуємо користувачів у документи та додаємо їх до колекції
	doc1, err := documentstore.MarshalDocument(user1)
	if err != nil {
		slog.Error("Error marshaling user1", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = usersCollection.Put(*doc1) // Розміновуємо doc1
	if err != nil {
		slog.Error("Error putting user1", slog.String("error", err.Error()))
		os.Exit(1)
	}

	doc2, err := documentstore.MarshalDocument(user2)
	if err != nil {
		slog.Error("Error marshaling user2", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = usersCollection.Put(*doc2) // Розміновуємо doc2
	if err != nil {
		slog.Error("Error putting user2", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Зберігаємо дамп сховища у файл
	dumpFile := "backup.json"
	err = store.DumpToFile(dumpFile)
	if err != nil {
		slog.Error("Error dumping to file", slog.String("error", err.Error()), slog.String("file", dumpFile))
		os.Exit(1)
	}
	slog.Info("Data dumped to file successfully", slog.String("file", dumpFile))

	// Відновлюємо Store з файлу
	restored, err := documentstore.NewStoreFromFile(dumpFile)
	if err != nil {
		slog.Error("Error restoring from file", slog.String("error", err.Error()), slog.String("file", dumpFile))
		os.Exit(1)
	}
	slog.Info("Data restored from file successfully", slog.String("file", dumpFile))

	// Отримуємо відновлену колекцію користувачів
	restoredUsersCollection, err := restored.GetCollection(usersCollectionName)
	if err != nil {
		slog.Error("Error getting restored users collection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Отримуємо всі документи з відновленої колекції користувачів
	restoredUsersDocs, err := restoredUsersCollection.GetAll()
	if err != nil {
		slog.Error("Error getting all restored users", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Перевіряємо, чи є користувачі у відновленій колекції
	if len(restoredUsersDocs) == 2 {
		slog.Info("Successfully restored users", slog.Int("count", len(restoredUsersDocs)))
		for _, doc := range restoredUsersDocs {
			var restoredUser User
			err = documentstore.UnmarshalDocument(&doc, &restoredUser)
			if err != nil {
				slog.Error("Error unmarshaling restored user", slog.String("error", err.Error()), slog.Any("document", doc))
				continue
			}
			fmt.Printf("Restored User: %+v\n", restoredUser)
		}
	} else {
		slog.Error("Failed to restore users", slog.Int("expected", 2), slog.Int("actual", len(restoredUsersDocs)))
	}

	fmt.Printf("Відновлено %d колекцій\n", restored.NumCollections())
}
