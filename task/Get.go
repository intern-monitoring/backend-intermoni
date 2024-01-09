package task

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTaskFromID(idtask primitive.ObjectID, db *mongo.Database) (task intermoni.Task, err error) {
	err = db.Collection("task").FindOne(context.Background(), bson.M{"_id": idtask}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return task, fmt.Errorf("task tidak ditemukan")
		}
		return task, fmt.Errorf("terjadi kesalahan")
	}
	return task, nil
}

func GetTaskByMahasiswa(iduser primitive.ObjectID, db *mongo.Database) (task intermoni.Task, err error) {
	mahasiswa,err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return task, fmt.Errorf("error GetTaskByMahasiswa get mahasiswa: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangByMahasiswa(mahasiswa.ID, db)
	if err != nil {
		return task, fmt.Errorf("error GetTaskByMahasiswa get mahasiswa magang: %s", err)
	}
	err = db.Collection("task").FindOne(context.Background(), bson.M{"mahasiswamagang._id": mahasiswa_magang.ID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return task, fmt.Errorf("task tidak ditemukan")
		}
		return task, fmt.Errorf("terjadi kesalahan")
	}
	return task, nil
}

func GetTaskByMahasiswaMagang(idmahasiswamagang primitive.ObjectID, db *mongo.Database) (task intermoni.Task, err error) {
	err = db.Collection("task").FindOne(context.Background(), bson.M{"mahasiswamagang._id": idmahasiswamagang}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return task, fmt.Errorf("task tidak ditemukan")
		}
		return task, fmt.Errorf("terjadi kesalahan")
	}
	return task, nil
}