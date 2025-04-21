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

func NewService(store *documentstore.Store, collectionName, primaryKey string) (*Service, error) {
	// Створюємо колекцію, якщо ще не існує
	_, err := store.CreateCollection(collectionName, &documentstore.CollectionConfig{
		PrimaryKey: primaryKey,
	})
	if err != nil && !errors.Is(err, documentstore.ErrCollectionAlreadyExists) {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	// Отримуємо колекцію
	coll, err := store.GetCollection(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &Service{coll: *coll}, nil
}
func (s *Service) CreateUser(name string) (*User, error) {
	user := &User{
		ID:   uuid.NewString(),
		Name: name,
	}

	// Серіалізуємо користувача у документ
	doc, err := documentstore.MarshalDocument(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	// Додаємо ключове поле вручну (оскільки MarshalDocument не знає про "key")
	doc.Fields["key"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: user.ID,
	}

	// Перевірка, чи документ успішно додано
	err = s.coll.Put(*doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	// Отримуємо всі документи з колекції
	documents, err := s.coll.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}

	users := make([]User, 0, len(documents)) // задаємо початкову ємність

	for _, doc := range documents {
		var user User
		err := documentstore.UnmarshalDocument(&doc, &user)
		if err != nil {
			continue // Пропускаємо, якщо не вдалося розпакувати
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

	var user User
	err = documentstore.UnmarshalDocument(doc, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	err := s.coll.Delete(userID)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}
