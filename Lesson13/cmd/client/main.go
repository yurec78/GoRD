package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"Lesson13/internal/documentstore"
	"Lesson13/internal/utils"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Не вдалося під'єднатися до сервера: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	fmt.Println("Клієнт під'єднано до сервера. Введіть команди (help для списку):")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // Видаляємо символ нового рядка

		if input == "exit" {
			break
		}

		var command utils.Command
		switch input {
		case "create_collection":
			var name, primaryKey string
			fmt.Print("Введіть назву колекції: ")
			fmt.Scanln(&name)
			fmt.Print("Введіть первинний ключ: ")
			fmt.Scanln(&primaryKey)
			command = utils.Command{
				Command: "create_collection",
				Payload: utils.CollectionConfigPayload{
					Name: name,
					Config: &documentstore.CollectionConfig{
						PrimaryKey: primaryKey,
					},
				},
			}
		case "delete_collection":
			var name string
			fmt.Print("Введіть назву колекції для видалення: ")
			fmt.Scanln(&name)
			command = utils.Command{
				Command: "delete_collection",
				Payload: utils.CollectionNamePayload{Name: name},
			}
		case "list_collections":
			command = utils.Command{Command: "list_collections", Payload: struct{}{}}
		case "put_document":
			var collectionName string
			fmt.Print("Введіть назву колекції: ")
			fmt.Scanln(&collectionName)
			fmt.Print("Введіть документ у форматі JSON: ")
			var doc map[string]interface{}
			decoder := json.NewDecoder(reader)
			if err := decoder.Decode(&doc); err != nil {
				fmt.Println("Невалідний формат JSON:", err)
				continue
			}
			command = utils.Command{
				Command: "put_document",
				Payload: utils.PutDocumentPayload{Collection: collectionName, Document: doc},
			}
		case "get_document":
			var collectionName, key string
			fmt.Print("Введіть назву колекції: ")
			fmt.Scanln(&collectionName)
			fmt.Print("Введіть ключ документа: ")
			fmt.Scanln(&key)
			command = utils.Command{
				Command: "get_document",
				Payload: utils.GetDeleteDocumentPayload{Collection: collectionName, Key: key},
			}
		case "delete_document":
			var collectionName, key string
			fmt.Print("Введіть назву колекції: ")
			fmt.Scanln(&collectionName)
			fmt.Print("Введіть ключ документа для видалення: ")
			fmt.Scanln(&key)
			command = utils.Command{
				Command: "delete_document",
				Payload: utils.GetDeleteDocumentPayload{Collection: collectionName, Key: key},
			}
		case "list_documents":
			var collectionName string
			fmt.Print("Введіть назву колекції: ")
			fmt.Scanln(&collectionName)
			command = utils.Command{
				Command: "list_documents",
				Payload: utils.CollectionNamePayload{Name: collectionName},
			}
		case "help":
			fmt.Println("Доступні команди:")
			fmt.Println("  create_collection - Створити нову колекцію")
			fmt.Println("  delete_collection - Видалити колекцію")
			fmt.Println("  list_collections  - Список усіх колекцій")
			fmt.Println("  put_document      - Додати/оновити документ")
			fmt.Println("  get_document      - Отримати документ за ключем")
			fmt.Println("  delete_document   - Видалити документ за ключем")
			fmt.Println("  list_documents    - Список усіх документів у колекції")
			fmt.Println("  exit              - Вийти з клієнта")
			continue // Пропускаємо відправку команди "help"
		default:
			fmt.Println("Невідома команда. Введіть 'help' для списку команд.")
			continue // Пропускаємо відправку невідомої команди
		}

		// Відправка команди на сервер
		jsonCommand, err := json.Marshal(command)
		if err != nil {
			fmt.Println("Помилка маршалінгу команди:", err)
			continue
		}
		_, err = conn.Write(append(jsonCommand, utils.Delimiter...))
		if err != nil {
			fmt.Println("Помилка відправки команди на сервер:", err)
			break
		}

		// Отримання відповіді від сервера
		responseStr, err := serverReader.ReadString(utils.Delimiter[0])
		if err != nil {
			fmt.Println("Помилка отримання відповіді від сервера:", err)
			break
		}
		responseStr = responseStr[:len(responseStr)-1] // Видаляємо роздільник

		var response utils.Response
		if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
			fmt.Println("Помилка розбору відповіді від сервера:", err)
			continue
		}

		// Виведення результату
		if response.Status == "ok" {
			if response.Result != nil {
				resultJSON, err := json.MarshalIndent(response.Result, "", "  ")
				if err != nil {
					fmt.Println("Помилка маршалінгу результату:", err)
					fmt.Printf("Результат: %+v\n", response.Result)
				} else {
					fmt.Println("Результат:")
					fmt.Println(string(resultJSON))
				}
			} else {
				fmt.Println("Успішно виконано.")
			}
		} else if response.Status == "error" {
			fmt.Println("Помилка від сервера:", response.Error.Message)
		}
	}

	fmt.Println("Клієнт завершив роботу.")
}
