package magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// by mahasiswa
func GetAllMagangByMahasiswa(db *mongo.Database) (magang []intermoni.Magang, err error) {
	mitra, err := GetMitraByMoU(db)
	if err != nil {
		return magang, err
	}
	for _, m := range mitra {
		mgn, err := GetAllMagangFromMitraBymahasiswa(m.ID, db)
		if err != nil {
			return magang, err
		}
		magang = append(magang, mgn...)
	}
	return magang, nil
}

func GetAllMagangFromMitraBymahasiswa(_id primitive.ObjectID, db *mongo.Database) (magang []intermoni.Magang, err error) {
	collection := db.Collection("magang")
	filter := bson.M{"mitra._id": _id}
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

func GetMitraByMoU(db *mongo.Database) (mitra []intermoni.Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{"mou": 1}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mitra, fmt.Errorf("error GetMitraByMoU mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mitra)
	if err != nil {
		return mitra, fmt.Errorf("error GetMitraByMoU context: %s", err)
	}
	return mitra, nil
}