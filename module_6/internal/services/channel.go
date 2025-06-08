package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"module_6/internal/models"
)

// ChannelService handles business logic for chat channels and messages
type ChannelService struct {
	messagesCollection *mongo.Collection
	channelsCollection *mongo.Collection
}

// NewChannelService creates a new ChannelService
func NewChannelService(messagesCol, channelsCol *mongo.Collection) *ChannelService {
	return &ChannelService{
		messagesCollection: messagesCol,
		channelsCollection: channelsCol,
	}
}

// GetMessages retrieves messages for a specific channel
func (s *ChannelService) GetMessages(ctx context.Context, channelID primitive.ObjectID, limit, offset int64) ([]models.Message, error) {
	if limit == 0 {
		limit = 100 // Default limit
	}
	if offset < 0 {
		offset = 0 // Default offset
	}

	findOptions := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Newest first

	cursor, err := s.messagesCollection.Find(ctx, bson.M{"channel_id": channelID}, findOptions)
	if err != nil {
		log.Printf("Error finding messages: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		log.Printf("Error decoding messages: %v", err)
		return nil, err
	}

	// Reverse the slice if you want oldest first for display, as we sorted by -1 (newest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// SendMessage saves a new message to the database
func (s *ChannelService) SendMessage(ctx context.Context, channelID, senderID primitive.ObjectID, senderUsername, content string) (*models.Message, error) {
	message := models.Message{
		ID:             primitive.NewObjectID(),
		ChannelID:      channelID,
		SenderID:       senderID,
		SenderUsername: senderUsername,
		Content:        content,
		Timestamp:      time.Now(),
	}

	_, err := s.messagesCollection.InsertOne(ctx, message)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return nil, err
	}
	log.Printf("Message sent to channel %s by %s", channelID.Hex(), senderUsername)
	return &message, nil
}

// GetChannelByName retrieves a channel by its name, or creates it if it doesn't exist
func (s *ChannelService) GetOrCreateChannelByName(ctx context.Context, channelName string) (*models.Channel, error) {
	var channel models.Channel
	err := s.channelsCollection.FindOne(ctx, bson.M{"name": channelName}).Decode(&channel)
	if err == nil {
		return &channel, nil // Channel found
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error finding channel by name: %v", err)
		return nil, err
	}

	// Channel not found, create a new one
	newChannel := models.Channel{
		ID:        primitive.NewObjectID(),
		Name:      channelName,
		CreatedAt: time.Now(),
	}
	_, err = s.channelsCollection.InsertOne(ctx, newChannel)
	if err != nil {
		log.Printf("Error creating new channel: %v", err)
		return nil, err
	}
	log.Printf("Created new channel: %s", newChannel.Name)
	return &newChannel, nil
}
