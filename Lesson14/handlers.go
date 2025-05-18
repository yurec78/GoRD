// handlers.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// writeJSONResponse надсилає JSON відповідь клієнту.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// makePostHandler загортає обробник, перевіряючи метод POST.
// Глобальна змінна 'store' буде визначена в main.go.
func makePostHandler(handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, StandardResponse{Ok: false, Error: "Only POST method is allowed"})
			return
		}
		handlerFunc(w, r)
	}
}

func handleCreateCollection(w http.ResponseWriter, r *http.Request) {
	var req CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name is required"})
		return
	}

	if err := store.CreateMongoCollection(r.Context(), req.CollectionName); err != nil {
		// Якщо помилка "collection already exists", можемо повернути 409 Conflict або інший код
		if strings.Contains(err.Error(), "already exists") {
			writeJSONResponse(w, http.StatusConflict, StandardResponse{Ok: false, Error: err.Error()})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		}
		return
	}
	writeJSONResponse(w, http.StatusCreated, StandardResponse{Ok: true}) // 201 Created
}

func handleListCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := store.ListMongoCollections(r.Context())
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ListCollectionsResponse{Ok: false, Error: err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, ListCollectionsResponse{Ok: true, Collections: collections})
}

func handleDeleteCollection(w http.ResponseWriter, r *http.Request) {
	var req DeleteCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name is required"})
		return
	}

	if err := store.DeleteMongoCollection(r.Context(), req.CollectionName); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, StandardResponse{Ok: true})
}

func handlePutDocument(w http.ResponseWriter, r *http.Request) {
	var req PutDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name is required"})
		return
	}
	if req.Document == nil { // bson.M може бути nil
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "document is required"})
		return
	}

	if err := store.PutMongoDocument(r.Context(), req.CollectionName, req.Document); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, StandardResponse{Ok: true})
}

func handleGetDocument(w http.ResponseWriter, r *http.Request) {
	var req GetDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" || req.DocumentID == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name and document_id are required"})
		return
	}

	doc, err := store.GetMongoDocument(r.Context(), req.CollectionName, req.DocumentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeJSONResponse(w, http.StatusNotFound, GetDocumentResponse{Ok: false, Error: err.Error()})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, GetDocumentResponse{Ok: false, Error: err.Error()})
		}
		return
	}
	writeJSONResponse(w, http.StatusOK, GetDocumentResponse{Ok: true, Document: doc})
}

func handleListDocuments(w http.ResponseWriter, r *http.Request) {
	var req ListDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name is required"})
		return
	}
	if req.Filter == nil {
		req.Filter = bson.M{} // Порожній фільтр, якщо не надано
	}

	docs, err := store.ListMongoDocuments(r.Context(), req.CollectionName, req.Filter, req.Limit, req.Skip)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ListDocumentsResponse{Ok: false, Error: err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, ListDocumentsResponse{Ok: true, Documents: docs})
}

func handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	var req DeleteDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" || req.DocumentID == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name and document_id are required"})
		return
	}

	if err := store.DeleteMongoDocument(r.Context(), req.CollectionName, req.DocumentID); err != nil {
		if strings.Contains(err.Error(), "not found for deletion") {
			writeJSONResponse(w, http.StatusNotFound, StandardResponse{Ok: false, Error: err.Error()})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		}
		return
	}
	writeJSONResponse(w, http.StatusOK, StandardResponse{Ok: true})
}

func handleCreateIndex(w http.ResponseWriter, r *http.Request) {
	var req CreateIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" || req.FieldName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name and field_name are required"})
		return
	}
	if req.Order == 0 {
		req.Order = 1 // За замовчуванням зростаючий
	}
	if req.Order != 1 && req.Order != -1 {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "order must be 1 (ascending) or -1 (descending)"})
		return
	}

	if err := store.CreateMongoIndex(r.Context(), req.CollectionName, req.FieldName, req.Unique, req.Order); err != nil {
		if strings.Contains(err.Error(), "already exist with different options") {
			writeJSONResponse(w, http.StatusConflict, StandardResponse{Ok: false, Error: err.Error()})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		}
		return
	}
	writeJSONResponse(w, http.StatusCreated, StandardResponse{Ok: true}) // 201 Created
}

func handleDeleteIndex(w http.ResponseWriter, r *http.Request) {
	var req DeleteIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "Invalid request body: " + err.Error()})
		return
	}
	if req.CollectionName == "" || req.IndexName == "" {
		writeJSONResponse(w, http.StatusBadRequest, StandardResponse{Ok: false, Error: "collection_name and index_name are required"})
		return
	}

	if err := store.DeleteMongoIndex(r.Context(), req.CollectionName, req.IndexName); err != nil {
		if strings.Contains(err.Error(), "index not found") {
			writeJSONResponse(w, http.StatusNotFound, StandardResponse{Ok: false, Error: err.Error()})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, StandardResponse{Ok: false, Error: err.Error()})
		}
		return
	}
	writeJSONResponse(w, http.StatusOK, StandardResponse{Ok: true})
}
