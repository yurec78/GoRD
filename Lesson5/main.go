package main

import (
	"Lesson5/documentstore"
	"Lesson5/users"
	"fmt"
	"log"
)

func main() {
	// Створюємо новий Store
	store := documentstore.NewStore()

	// Ініціалізуємо сервіс користувачів (він сам створить колекцію)
	userService, err := users.NewService(store, "users", "key")
	if err != nil {
		log.Fatalf("failed to initialize user service: %v", err)
	}

	// Створюємо нового користувача
	user1, err := userService.CreateUser("John Doe")
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Printf("Created user: %v\n", user1)

	// Створюємо ще одного користувача
	user2, err := userService.CreateUser("Jane Smith")
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Printf("Created user: %v\n", user2)

	// Отримуємо список усіх користувачів
	usersList, err := userService.ListUsers()
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}
	fmt.Println("List of users:")
	for _, user := range usersList {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

	// Отримуємо одного користувача за його ID
	userID := user1.ID
	user, err := userService.GetUser(userID)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}
	fmt.Printf("Fetched user by ID: ID: %s, Name: %s\n", user.ID, user.Name)

	// Видаляємо користувача
	err = userService.DeleteUser(userID)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
	fmt.Printf("Deleted user with ID: %s\n", userID)

	// Перевіряємо, чи користувач залишився після видалення
	_, err = userService.GetUser(userID)
	if err != nil {
		fmt.Printf("User with ID %s not found (as expected)\n", userID)
	}
}
