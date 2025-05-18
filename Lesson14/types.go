package main

import "go.mongodb.org/mongo-driver/bson"

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type DeleteCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type PutDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	Document       bson.M `json:"document"`
}

type GetDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	DocumentID     string `json:"document_id"`
}

type ListDocumentsRequest struct {
	CollectionName string `json:"collection_name"`
	Filter         bson.M `json:"filter,omitempty"`
	Limit          int64  `json:"limit,omitempty"`
	Skip           int64  `json:"skip,omitempty"`
}

type DeleteDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	DocumentID     string `json:"document_id"`
}

type CreateIndexRequest struct {
	CollectionName string `json:"collection_name"`
	FieldName      string `json:"field_name"`
	Unique         bool   `json:"unique,omitempty"`
	Order          int    `json:"order,omitempty"` // 1 for asc, -1 for desc. Default to 1 (asc)
}

type DeleteIndexRequest struct {
	CollectionName string `json:"collection_name"`
	IndexName      string `json:"index_name"`
}

// --- Response Structs ---

type StandardResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type ListCollectionsResponse struct {
	Ok          bool     `json:"ok"`
	Collections []string `json:"collections,omitempty"`
	Error       string   `json:"error,omitempty"`
}

type GetDocumentResponse struct {
	Ok       bool   `json:"ok"`
	Document bson.M `json:"document,omitempty"`
	Error    string `json:"error,omitempty"`
}

type ListDocumentsResponse struct {
	Ok        bool     `json:"ok"`
	Documents []bson.M `json:"documents,omitempty"`
	Error     string   `json:"error,omitempty"`
}
