package task

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditTaskOlehMentor(idmahasiswamagang primitive.ObjectID, db *mongo.Database, updatedDoc intermoni.Task) error {
	if updatedDoc.Tasks == nil {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	task, err := GetTaskByIDMahasiswaMagang(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": idmahasiswamagang,
		},
		"tasks": updatedDoc.Tasks,
	}
	err = intermoni.UpdateOneDoc(task.ID, db, "task", data)
	if err != nil {
		return err
	}
	return nil
}

func GetTaskByIDMahasiswaMagang(idmahasiswamagang primitive.ObjectID, db *mongo.Database) (task intermoni.Task, err error) {
	err = db.Collection("task").FindOne(context.Background(), bson.M{"mahasiswamagang._id": idmahasiswamagang}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return task, fmt.Errorf("task tidak ditemukan")
		}
		return task, fmt.Errorf("terjadi kesalahan")
	}
	return task, nil
}