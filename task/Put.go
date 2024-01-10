package task

import (
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
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": idmahasiswamagang,
		},
		"tasks": updatedDoc.Tasks,
	}
	err := intermoni.UpdateOneDoc(idmahasiswamagang, db, "task", data)
	if err != nil {
		return err
	}
	return nil
}