package storage

import (
	"context"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Delete all challenges related to this section first
	challengeCollection := s.db.MongoDB.GetCollection("challenges")
	challengeFilter := bson.M{"section_id": id}
	_, err := challengeCollection.DeleteMany(ctx, challengeFilter)
	if err != nil {
		return err
	}

	// Then delete the section
	collection := s.db.MongoDB.GetCollection(model.SectionName)
	filter := bson.M{"_id": id}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
