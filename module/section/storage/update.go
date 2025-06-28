package storage

import (
	"context"
	"hub-service/module/section/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Update(ctx context.Context, id primitive.ObjectID, data *model.SectionUpdate) error {
	now := time.Now()
	data.UpdatedAt = &now

	collection := s.db.MongoDB.GetCollection(model.SectionName)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": data}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
