package pembimbing

import (
	"context"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// by admin
func GetAllPembimbingByAdmin(db *mongo.Database) (pembimbing []intermoni.Pembimbing, err error) {
	filter := bson.M{}
	cursor, err := db.Collection("pembimbing").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &pembimbing); err != nil {
		return nil, err
	}
	return pembimbing, nil
}