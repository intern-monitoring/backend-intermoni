package task

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TambahTaskOlehMentor(idmahasiswamagang primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Task) error {
	if insertedDoc.Tasks == nil {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if CheckTaskByMahasiswaMagang(idmahasiswamagang, db) {
		return fmt.Errorf("mahasiswa sudah memiliki task")
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": idmahasiswamagang,
		},
		"tasks": insertedDoc.Tasks,
	}
	_, err := intermoni.InsertOneDoc(db, "task", data)
	if err != nil {
		return err
	}
	return nil
}

func CheckTaskByMahasiswaMagang(idmahasiswamagang primitive.ObjectID, db *mongo.Database) bool {
	filter := bson.M{"mahasiswamagang._id": idmahasiswamagang}
	err := db.Collection("task").FindOne(context.Background(), filter).Decode(&task)
	return err == nil
}