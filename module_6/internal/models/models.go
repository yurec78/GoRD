package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Password  string             `bson:"password" json:"-" validate:"required,min=6"` // "-" щоб не серіалізувати пароль
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID      primitive.ObjectID `bson:"channel_id" json:"channel_id" validate:"required"`
	SenderID       primitive.ObjectID `bson:"sender_id" json:"sender_id" validate:"required"`
	SenderUsername string             `bson:"sender_username" json:"sender_username"`
	Content        string             `bson:"content" json:"content" validate:"required,min=1,max=500"`
	Timestamp      time.Time          `bson:"timestamp" json:"timestamp"`
}

type Channel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required,min=3,max=50"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SendMessageRequest struct {
	ChannelID string `json:"channel_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

type GetHistoryRequest struct {
	ChannelID string `query:"channel_id" validate:"required"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
}

type WebSocketMessage struct {
	ChannelID      string    `json:"channel_id"`
	Content        string    `json:"content"`
	SenderID       string    `json:"sender_id,omitempty"`
	SenderUsername string    `json:"sender_username,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}
