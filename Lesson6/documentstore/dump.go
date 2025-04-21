package documentstore

import (
	"encoding/json"
	"log"
	"os"
)

// Dump повертає дамп (JSON) усього Store: колекцій та документів
func (s *Store) Dump() ([]byte, error) {
	dump, err := json.Marshal(s)
	if err != nil {
		log.Printf("STORE DUMP FAILED: Помилка маршалінгу JSON: %v", err)
		return nil, err
	}
	log.Println("STORE DUMPED: Створено дамп сховища")
	return dump, nil
}

// DumpToFile зберігає дамп Store у файл
func (s *Store) DumpToFile(filename string) error {
	data, err := s.Dump()
	if err != nil {
		log.Printf("STORE DUMP TO FILE FAILED: Помилка отримання дампу: %v", err)
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("STORE DUMP TO FILE FAILED: Помилка запису у файл '%s': %v", filename, err)
		return err
	}
	log.Printf("STORE DUMPED TO FILE: Дамп сховища збережено у файл '%s'", filename)
	return nil
}

// NewStoreFromDump створює новий Store із JSON-дампу
func NewStoreFromDump(dump []byte) (*Store, error) {
	var store Store
	err := json.Unmarshal(dump, &store)
	if err != nil {
		log.Printf("STORE RESTORE FAILED: Помилка демаршалінгу JSON: %v", err)
		return nil, err
	}
	log.Println("STORE RESTORED FROM DUMP: Сховище успішно відновлено з дампу")
	return &store, nil
}

// NewStoreFromFile читає дамп із файлу та відновлює Store
func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("STORE RESTORE FROM FILE FAILED: Помилка читання файлу '%s': %v", filename, err)
		return nil, err
	}
	store, err := NewStoreFromDump(data)
	if err != nil {
		return nil, err
	}
	log.Printf("STORE RESTORED FROM FILE: Сховище успішно відновлено з файлу '%s'", filename)
	return store, nil
}
