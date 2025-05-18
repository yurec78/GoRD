package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"Lesson11/documentstore"
)

func main() {
	store := documentstore.NewStore()
	config := &documentstore.CollectionConfig{
		PrimaryKey: "ID",
	}
	err := store.CreateCollection("mycollection", config)
	if err != nil {
		log.Fatalf("Не вдалося створити колекцію 'mycollection': %v", err)
	}
	collection, err := store.GetCollection("mycollection")
	if err != nil {
		log.Fatalf("Не вдалося отримати колекцію 'mycollection': %v", err)
	}

	fmt.Println("Починаємо запуск горутин...")
	var wg sync.WaitGroup
	numGoroutines := 1000

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Printf("Горутина %d запущена", id)

			key := fmt.Sprintf("doc-%d", rand.Intn(50))
			doc := documentstore.Document{
				Fields: map[string]documentstore.DocumentField{
					"ID":    {Type: documentstore.DocumentFieldTypeString, Value: key},
					"Value": {Type: documentstore.DocumentFieldTypeString, Value: fmt.Sprintf("value-%d", rand.Intn(100))},
				},
			}

			op := rand.Intn(3)
			switch op {
			case 0:
				err := collection.Put(doc)
				if err != nil {
					log.Printf("Горутина %d: Помилка Put: %v", id, err)
				} else {
					log.Printf("Горутина %d: Додано документ з ключем %s", id, key)
				}
			case 1:
				_, err := collection.Get(key)
				if err != nil && !errors.Is(err, documentstore.ErrDocumentNotFound) {
					log.Printf("Горутина %d: Помилка Get: %v", id, err)
				} else if err == nil {
					log.Printf("Горутина %d: Отримано документ з ключем %s", id, key)
				}
			case 2:
				err := collection.Delete(key)
				if err != nil && !errors.Is(err, documentstore.ErrDocumentNotFound) {
					log.Printf("Горутина %d: Помилка Delete: %v", id, err)
				} else if err == nil {
					log.Printf("Горутина %d: Видалено документ з ключем %s", id, key)
				}
			}

			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			log.Printf("Горутина %d завершила виконання", id)
		}(i)
	}

	fmt.Printf("Кількість горутин перед очікуванням: %d\n", runtime.NumGoroutine())
	wg.Wait()
	fmt.Println("Усі горутини завершили виконання.")
	fmt.Printf("Кількість запущених горутин після очікування: %d\n", runtime.NumGoroutine())
}
