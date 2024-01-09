package task

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditTaskOlehMentor(idtask primitive.ObjectID, db *mongo.Database, updatedDoc intermoni.Task) error {
	if updatedDoc.Tasks == nil {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	task, err := GetTaskFromID(idtask, db)
	if err != nil {
		return err
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": task.MahasiswaMagang.ID,
		},
		"tasks": updatedDoc.Tasks,
	}
	err = intermoni.UpdateOneDoc(idtask, db, "task", data)
	if err != nil {
		return err
	}
	return nil
}