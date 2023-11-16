package magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mahasiswa
func GetAllMagangByMahasiswa(db *mongo.Database) (magang []intermoni.Magang, err error) {
	collection := db.Collection("magang")
	magang, err = GetAllMagang(db)
	if err != nil {
		return magang, err
	}
	filter := bson.M{"mitra.mou": 1}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagangByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &magang)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagangByMahasiswa context: %s", err)
	}
	return magang, nil
}

// by mitra
func GetAllMagangByMitra(_id primitive.ObjectID, db *mongo.Database) (magang []intermoni.Magang, err error) {
	collection := db.Collection("magang")
	mitra, err := intermoni.GetMitraFromAkun(_id, db)
	if err != nil {
		return magang, err
	}
	filter := bson.M{"mitra._id": mitra.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return magang, fmt.Errorf("error GetMagangByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &magang)
	if err != nil {
		return magang, fmt.Errorf("error GetMagangByMitra context: %s", err)
	}
	return magang, nil
}

// by admin
func GetAllMagang(db *mongo.Database) (magang []intermoni.Magang, err error) {
	collection := db.Collection("magang")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagang mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &magang)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagang context: %s", err)
	}
	for _, m := range magang {
		mitra, err := intermoni.GetMitraFromID(m.Mitra.ID, db)
		if err != nil {
			fmt.Println(m.Mitra.ID)
			return magang, fmt.Errorf("error GetAllMagang get mitra: %s", err)
		}
		m.Mitra = mitra
		magang = append(magang, m)
		magang = magang[1:]
	}
	return magang, nil
}