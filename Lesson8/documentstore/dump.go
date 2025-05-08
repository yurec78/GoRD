package documentstore

import (
	"encoding/json"
	"log/slog"
	"os"
)

// Dump повертає дамп (JSON) усього Store: колекцій та документів
func (s *Store) Dump() ([]byte, error) {
	dump, err := json.Marshal(s)
	if err != nil {
		slog.Error("STORE DUMP FAILED", slog.Any("error", err), slog.String("message", "Помилка маршалінгу JSON"))
		return nil, err
	}
	slog.Debug("STORE DUMPED", slog.String("message", "Створено дамп сховища"), slog.String("dump_json", string(dump)))
	return dump, nil
}

// DumpToFile зберігає дамп Store у файл
func (s *Store) DumpToFile(filename string) error {
	data, err := s.Dump()
	if err != nil {
		slog.Error("STORE DUMP TO FILE FAILED", slog.Any("error", err), slog.String("message", "Помилка отримання дампу"))
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		slog.Error("STORE DUMP TO FILE FAILED", slog.String("filename", filename), slog.Any("error", err), slog.String("message", "Помилка запису у файл"))
		return err
	}
	slog.Info("STORE DUMPED TO FILE", slog.String("filename", filename), slog.String("message", "Дамп сховища збережено у файл"))
	return nil
}

// NewStoreFromDump створює новий Store із JSON-дампу
func NewStoreFromDump(dump []byte) (*Store, error) {
	var store Store
	err := json.Unmarshal(dump, &store)
	if err != nil {
		slog.Error("STORE RESTORE FAILED", slog.Any("error", err), slog.String("dump_json", string(dump)), slog.String("message", "Помилка демаршалінгу JSON"))
		return nil, err
	}
	slog.Info("STORE RESTORED FROM DUMP", slog.String("message", "Сховище успішно відновлено з дампу"))
	slog.Debug("STORE RESTORED FROM DUMP", slog.Any("restored_store", store))
	return &store, nil
}

// NewStoreFromFile читає дамп із файлу та відновлює Store
func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("STORE RESTORE FROM FILE FAILED", slog.String("filename", filename), slog.Any("error", err), slog.String("message", "Помилка читання файлу"))
		return nil, err
	}
	slog.Debug("STORE RESTORE FROM FILE", slog.String("filename", filename), slog.String("file_content", string(data)), slog.String("message", "Вміст файлу прочитано"))
	store, err := NewStoreFromDump(data)
	if err != nil {
		return nil, err
	}
	slog.Info("STORE RESTORED FROM FILE", slog.String("filename", filename), slog.String("message", "Сховище успішно відновлено з файлу"))
	return store, nil
}
