package users

import (
	"fmt"
	"log/slog"

	"Lesson6/documentstore"
	"errors"
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
	coll *documentstore.Collection // Змінено на вказівник для уникнення копіювання
}

func NewService(store *documentstore.Store, collectionName, primaryKey string) (*Service, error) {
	// Створюємо колекцію, якщо ще не існує
	err := store.CreateCollection(collectionName, &documentstore.CollectionConfig{
		PrimaryKey: primaryKey,
	})
	if err != nil && !errors.Is(err, documentstore.ErrCollectionAlreadyExists) {
		slog.Error("SERVICE CREATION FAILED", slog.String("collection", collectionName), slog.Any("error", err), slog.String("message", fmt.Sprintf("Не вдалося створити колекцію '%s'", collectionName)))
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}
	if err == nil {
		slog.Info("SERVICE CREATION", slog.String("collection", collectionName), slog.String("message", fmt.Sprintf("Колекцію '%s' створено", collectionName)))
	} else if errors.Is(err, documentstore.ErrCollectionAlreadyExists) {
		slog.Info("SERVICE CREATION", slog.String("collection", collectionName), slog.String("message", fmt.Sprintf("Колекція '%s' вже існує", collectionName)))
	}

	// Отримуємо колекцію
	coll, err := store.GetCollection(collectionName)
	if err != nil {
		slog.Error("SERVICE CREATION FAILED", slog.String("collection", collectionName), slog.Any("error", err), slog.String("message", fmt.Sprintf("Не вдалося отримати колекцію '%s'", collectionName)))
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}
	slog.Debug("SERVICE CREATION", slog.String("collection", collectionName), slog.String("message", fmt.Sprintf("Колекцію '%s' отримано", collectionName)), slog.Any("collection_object", coll))

	slog.Info("SERVICE CREATED", slog.String("collection", collectionName), slog.String("message", fmt.Sprintf("Створено новий сервіс для колекції '%s'", collectionName)))
	return &Service{coll: coll}, nil
}

func (s *Service) CreateUser(name string) (*User, error) {
	user := &User{
		ID:   uuid.NewString(),
		Name: name,
	}
	slog.Debug("CREATE USER", slog.Any("user", user), slog.String("message", fmt.Sprintf("Створення користувача з ID '%s' та ім'ям '%s'", user.ID, user.Name)))

	// Серіалізуємо користувача у документ
	doc, err := documentstore.MarshalDocument(user)
	if err != nil {
		slog.Error("CREATE USER FAILED", slog.Any("user", user), slog.Any("error", err), slog.String("message", "Не вдалося маршалізувати користувача"))
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}
	slog.Debug("CREATE USER", slog.Any("user", user), slog.Any("document", doc), slog.String("message", "Користувача успішно маршалізовано"))

	// Додаємо ключове поле вручну (оскільки MarshalDocument не знає про "key")
	doc.Fields["key"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: user.ID,
	}
	slog.Debug("CREATE USER", slog.Any("user", user), slog.Any("document_with_key", doc), slog.String("message", "Додано ключове поле до документа"))

	// Перевірка, чи документ успішно додано
	err = s.coll.Put(*doc)
	if err != nil {
		slog.Error("CREATE USER FAILED", slog.Any("user", user), slog.Any("document", doc), slog.Any("error", err), slog.String("message", "Не вдалося додати користувача до колекції"))
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	slog.Info("USER CREATED", slog.Any("user", user), slog.String("message", fmt.Sprintf("Користувача '%s' з ID '%s' створено", user.Name, user.ID)))

	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	slog.Debug("LIST USERS", slog.String("collection", s.coll.Name()), slog.String("message", "Отримання списку всіх користувачів"))
	// Отримуємо всі документи з колекції
	documents, err := s.coll.GetAll()
	if err != nil {
		slog.Error("LIST USERS FAILED", slog.String("collection", s.coll.Name()), slog.Any("error", err), slog.String("message", "Не вдалося отримати документи з колекції"))
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}
	slog.Debug("LIST USERS", slog.String("collection", s.coll.Name()), slog.Int("count", len(documents)), slog.String("message", fmt.Sprintf("Отримано %d документів", len(documents))))

	users := make([]User, 0, len(documents)) // задаємо початкову ємність

	for _, doc := range documents {
		var user User
		err := documentstore.UnmarshalDocument(&doc, &user)
		if err != nil {
			slog.Warn("LIST USERS", slog.String("collection", s.coll.Name()), slog.Any("document", doc), slog.Any("error", err), slog.String("message", "Не вдалося демаршалізувати документ користувача, пропущено"))
			continue // Пропускаємо, якщо не вдалося розпакувати
		}
		users = append(users, user)
		slog.Debug("LIST USERS", slog.String("collection", s.coll.Name()), slog.Any("user", user), slog.String("message", "Користувача демаршалізовано"))
	}
	slog.Info("LIST USERS", slog.String("collection", s.coll.Name()), slog.Int("count", len(users)), slog.String("message", fmt.Sprintf("Отримано %d користувачів", len(users))))
	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	slog.Debug("GET USER", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.String("message", fmt.Sprintf("Отримання користувача з ID '%s'", userID)))
	// Отримуємо документ за userID
	doc, err := s.coll.Get(userID)
	if err != nil {
		slog.Warn("GET USER FAILED", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.Any("error", err), slog.String("message", fmt.Sprintf("Користувача з ID '%s' не знайдено", userID)))
		return nil, ErrUserNotFound
	}
	slog.Debug("GET USER", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.Any("document", doc), slog.String("message", "Документ користувача отримано"))

	var user User
	err = documentstore.UnmarshalDocument(doc, &user)
	if err != nil {
		slog.Error("GET USER FAILED", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.Any("document", doc), slog.Any("error", err), slog.String("message", "Не вдалося демаршалізувати користувача"))
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	slog.Info("USER FOUND", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.Any("user", user), slog.String("message", fmt.Sprintf("Користувача з ID '%s' знайдено", userID)))
	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	slog.Debug("DELETE USER", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.String("message", fmt.Sprintf("Видалення користувача з ID '%s'", userID)))
	err := s.coll.Delete(userID)
	if err != nil {
		slog.Warn("DELETE USER FAILED", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.Any("error", err), slog.String("message", fmt.Sprintf("Не вдалося видалити користувача з ID '%s'", userID)))
		return ErrUserNotFound
	}
	slog.Info("USER DELETED", slog.String("collection", s.coll.Name()), slog.String("userID", userID), slog.String("message", fmt.Sprintf("Користувача з ID '%s' видалено", userID)))
	return nil
}
