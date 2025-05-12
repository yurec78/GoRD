package utils

import "Lesson12/internal/documentstore" // Замініть на ваш фактичний шлях

// Command - структура для команд від клієнта
type Command struct {
	Command string      `json:"command"`
	Payload interface{} `json:"payload"`
}

// Response - структура для відповідей від сервера
type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}

// Error - структура для повідомлень про помилки
type Error struct {
	Message string `json:"message"`
}

// CollectionConfigPayload - структура для payload команди create_collection
type CollectionConfigPayload struct {
	Name   string                          `json:"name"`
	Config *documentstore.CollectionConfig `json:"config"`
}

// CollectionNamePayload - структура для payload команд delete_collection, list_documents
type CollectionNamePayload struct {
	Name string `json:"name"`
}

// PutDocumentPayload - структура для payload команди put_document
type PutDocumentPayload struct {
	Collection string                 `json:"collection"`
	Document   map[string]interface{} `json:"document"`
}

// GetDeleteDocumentPayload - структура для payload команд get_document, delete_document
type GetDeleteDocumentPayload struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}

// ListCollectionsResult - структура для результату команди list_collections
type ListCollectionsResult struct {
	Collections []string `json:"collections"`
}

// GetDocumentResult - структура для результату команди get_document
type GetDocumentResult struct {
	Document map[string]interface{} `json:"document"`
}

// ListDocumentsResult - структура для результату команди list_documents
type ListDocumentsResult struct {
	Documents []map[string]interface{} `json:"documents"`
}

// GenericResult - структура для простих результатів (ok/error повідомлення)
type GenericResult struct {
	Message string `json:"message"`
}
