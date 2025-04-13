package users

import (
	"Lesson5/documentstore"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidDocumentField = errors.New("invalid document field")
	ErrCollectionNotFound   = errors.New("collection not found")
	ErrDocumentNotFound     = errors.New("document not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll documentstore.Collection
}

func NewService(coll documentstore.Collection) *Service {
	return &Service{coll: coll}
}

func (s *Service) CreateUser(name string) (*User, error) {
	user := &User{
		ID:   uuid.NewString(),
		Name: name,
	}

	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: user.ID,
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: user.Name,
			},
		},
	}

	// Перевірка, чи документ успішно додано
	err := s.coll.Put(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	var users []User

	// Отримуємо всі документи з колекції
	documents, err := s.coll.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}

	for _, doc := range documents {
		// Перевірка наявності поля "name"
		nameField, ok := doc.Fields["name"]
		if !ok || nameField.Type != documentstore.DocumentFieldTypeString {
			continue // Якщо поля немає або воно не типу "string", пропускаємо документ
		}

		name, ok := nameField.Value.(string)
		if !ok {
			continue // Якщо значення "name" не рядок, пропускаємо
		}

		// Додаємо користувача в список
		user := User{
			ID:   doc.Fields["key"].Value.(string), // Поле "key" містить ID
			Name: name,
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	// Отримуємо документ за userID
	doc, err := s.coll.Get(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Перевіряємо наявність полів "id" і "name"
	idField, ok := doc.Fields["id"]
	if !ok || idField.Type != documentstore.DocumentFieldTypeString {
		return nil, fmt.Errorf("missing or invalid 'id' field")
	}

	nameField, ok := doc.Fields["name"]
	if !ok || nameField.Type != documentstore.DocumentFieldTypeString {
		return nil, fmt.Errorf("missing or invalid 'name' field")
	}

	user := &User{
		ID:   idField.Value.(string),
		Name: nameField.Value.(string),
	}

	return user, nil
}

func (s *Service) DeleteUser(userID string) error {
	// Перевірка наявності документа
	doc, ok := s.coll.Documents[userID]
	if !ok {
		return ErrUserNotFound
	}

	// Перевірка поля "key"
	keyField, ok := doc.Fields["key"]
	if !ok || keyField.Type != documentstore.DocumentFieldTypeString {
		return fmt.Errorf("invalid or missing 'key' field")
	}

	// Видалення документа
	delete(s.coll.Documents, userID)

	return nil
}
