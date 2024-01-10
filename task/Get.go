package task

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTaskFromID(idtask primitive.ObjectID, db *mongo.Database) (bson.M, error) {
	var task intermoni.Task
	err := db.Collection("task").FindOne(context.Background(), bson.M{"_id": idtask}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bson.M{}, fmt.Errorf("task tidak ditemukan")
		}
		return bson.M{}, fmt.Errorf("terjadi kesalahan")
	}
	data := bson.M{
		"_id" : task.ID,
		"mahasiswamagang" : bson.M{
			"_id" : task.MahasiswaMagang.ID,
		},
		"tasks" : task.Tasks,
	}
	return data, nil
}

func GetTaskByMahasiswa(iduser primitive.ObjectID, db *mongo.Database) (bson.M, error) {
	var task intermoni.Task
	mahasiswa,err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return bson.M{}, fmt.Errorf("error GetTaskByMahasiswa get mahasiswa: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangByMahasiswa(mahasiswa.ID, db)
	if err != nil {
		return bson.M{}, fmt.Errorf("error GetTaskByMahasiswa get mahasiswa magang: %s", err)
	}
	err = db.Collection("task").FindOne(context.Background(), bson.M{"mahasiswamagang._id": mahasiswa_magang.ID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bson.M{}, fmt.Errorf("task tidak ditemukan")
		}
		return bson.M{}, fmt.Errorf("terjadi kesalahan")
	}
	data := bson.M{
		"_id" : task.ID,
		"mahasiswamagang" : bson.M{
			"_id" : task.MahasiswaMagang.ID,
		},
		"tasks" : task.Tasks,
	}
	return data, nil
}

func GetTaskByMahasiswaMagang(idmahasiswamagang primitive.ObjectID, db *mongo.Database) (bson.M, error) {
	var task intermoni.Task
	err := db.Collection("task").FindOne(context.Background(), bson.M{"mahasiswamagang._id": idmahasiswamagang}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bson.M{}, fmt.Errorf("task tidak ditemukan")
		}
		return bson.M{}, fmt.Errorf("terjadi kesalahan")
	}
	data := bson.M{
		"_id" : task.ID,
		"mahasiswamagang" : bson.M{
			"_id" : task.MahasiswaMagang.ID,
		},
		"tasks" : task.Tasks,
	}
	return data, nil
}