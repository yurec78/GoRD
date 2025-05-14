package main

import (
	"Lesson12/internal/documentstore"
	"Lesson12/internal/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn, store *documentstore.Store) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString(utils.Delimiter[0])
		if err != nil {
			log.Printf("Клієнт %s від'єднався: %v", conn.RemoteAddr(), err)
			return
		}
		message = message[:len(message)-1] // Видаляємо роздільник

		var command utils.Command
		if err := json.Unmarshal([]byte(message), &command); err != nil {
			log.Printf("Помилка розбору команди від %s: %v", conn.RemoteAddr(), err)
			response := utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідна команда"}}
			jsonResponse, _ := json.Marshal(response)
			conn.Write(append(jsonResponse, utils.Delimiter...))
			continue
		}

		response := processCommand(command, store)
		jsonResponse, _ := json.Marshal(response)
		conn.Write(append(jsonResponse, utils.Delimiter...))
	}
}

func processCommand(command utils.Command, store *documentstore.Store) utils.Response {
	switch command.Command {
	case "create_collection":
		var payload utils.CollectionConfigPayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для create_collection"}}
		}
		err := store.CreateCollection(payload.Name, payload.Config)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		return utils.Response{Status: "ok", Result: &utils.GenericResult{Message: fmt.Sprintf("Колекцію '%s' створено успішно", payload.Name)}}

	case "delete_collection":
		var payload utils.CollectionNamePayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для delete_collection"}}
		}
		err := store.DeleteCollection(payload.Name)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		return utils.Response{Status: "ok", Result: &utils.GenericResult{Message: fmt.Sprintf("Колекцію '%s' видалено", payload.Name)}}

	case "list_collections":
		collections := store.GetAllCollections()
		names := make([]string, 0, len(collections))
		for name := range collections {
			names = append(names, name)
		}
		return utils.Response{Status: "ok", Result: &utils.ListCollectionsResult{Collections: names}}

	case "put_document":
		var payload utils.PutDocumentPayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для put_document"}}
		}
		collection, err := store.GetCollection(payload.Collection)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		doc, err := documentstore.MarshalDocument(payload.Document)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		err = collection.Put(*doc)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		return utils.Response{Status: "ok", Result: &utils.GenericResult{Message: fmt.Sprintf("Документ додано/оновлено в колекції '%s'", payload.Collection)}}

	case "get_document":
		var payload utils.GetDeleteDocumentPayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для get_document"}}
		}
		collection, err := store.GetCollection(payload.Collection)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		docPtr, err := collection.Get(payload.Key)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		var result map[string]interface{}
		if err := documentstore.UnmarshalDocument(docPtr, &result); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		return utils.Response{Status: "ok", Result: &utils.GetDocumentResult{Document: result}}

	case "delete_document":
		var payload utils.GetDeleteDocumentPayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для delete_document"}}
		}
		collection, err := store.GetCollection(payload.Collection)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		err = collection.Delete(payload.Key)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		return utils.Response{Status: "ok", Result: &utils.GenericResult{Message: fmt.Sprintf("Документ з ключем '%s' видалено з колекції '%s'", payload.Key, payload.Collection)}}

	case "list_documents":
		var payload utils.CollectionNamePayload
		if err := json.Unmarshal([]byte(command.Payload.(string)), &payload); err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: "Невалідний payload для list_documents"}}
		}
		collection, err := store.GetCollection(payload.Name)
		if err != nil {
			return utils.Response{Status: "error", Error: &utils.Error{Message: err.Error()}}
		}
		docs := collection.List()
		results := make([]map[string]interface{}, 0, len(docs))
		for _, doc := range docs {
			var result map[string]interface{}
			if err := documentstore.UnmarshalDocument(&doc, &result); err != nil {
				log.Printf("Помилка демаршалінгу документа: %v", err)
				continue
			}
			results = append(results, result)
		}
		return utils.Response{Status: "ok", Result: &utils.ListDocumentsResult{Documents: results}}

	default:
		return utils.Response{Status: "error", Error: &utils.Error{Message: fmt.Sprintf("Невідома команда: %s", command.Command)}}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Не вдалося запустити сервер: %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	store := documentstore.NewStore()
	log.Println("Сервер запущено та слухає на :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Помилка при прийнятті з'єднання: %v", err)
			continue
		}
		go handleConnection(conn, store)
	}
}
