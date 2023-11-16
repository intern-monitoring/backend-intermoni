package mitra

import (
	"context"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by admin
func GetMitraFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (mitra intermoni.Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&mitra)
	if err != nil {
		return mitra, err
	}
	user, err := intermoni.GetUserFromID(mitra.Akun.ID, db)
	if err != nil {
		return mitra, err
	}
	akun := intermoni.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	mitra.Akun = akun
	return mitra, nil
}

func GetAllMitraByAdmin(db *mongo.Database) (mitra []intermoni.Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mitra, err
	}
	err = cursor.All(context.Background(), &mitra)
	if err != nil {
		return mitra, err
	}
	return mitra, nil
}