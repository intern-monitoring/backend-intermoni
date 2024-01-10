package report

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mahasiswa
func GetAllReport(_id primitive.ObjectID, db *mongo.Database) (data []bson.M, err error) {
	var report []intermoni.Report
	collection := db.Collection("report")
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(_id, db)
	if err != nil {
		return data, fmt.Errorf("error GetAllReport get mahasiswa: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangByMahasiswa(mahasiswa.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetAllReport get mahasiswa magang: %s", err)
	}
	filter := bson.M{"mahasiswamagang._id": mahasiswa_magang.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa context: %s", err)
	}
	for _, r := range report {
		mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(r.MahasiswaMagang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa magang: %s", err)
		}
		magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get magang: %s", err)
		}
		datareport := bson.M{
			"magang": magang,
			"_id": r.ID,
			"task": r.Task,
			"deskripsi": r.Deskripsi,
			"hasil": r.Hasil,
			"kehadiran": r.Kehadiran,
			"createdat": r.CreatedAt,
			"updatedat": r.UpdatedAt,
		}
		data = append(data, datareport)
	}
	return data, nil
}

func GetAllReportOlehPembimbing(_id primitive.ObjectID, db *mongo.Database) (data []bson.M, err error) {
	var report []intermoni.Report
	collection := db.Collection("report")
	filter := bson.M{"mahasiswamagang._id": _id}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa context: %s", err)
	}
	for _, r := range report {
		mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(r.MahasiswaMagang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa magang: %s", err)
		}
		magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get magang: %s", err)
		}
		datareport := bson.M{
			"magang": magang,
			"_id": r.ID,
			"task": r.Task,
			"deskripsi": r.Deskripsi,
			"kehadiran": r.Kehadiran,
			"createdat": r.CreatedAt,
			"updatedat": r.UpdatedAt,
		}
		data = append(data, datareport)
	}
	return data, nil
}

func GetAllReportOlehMentor(_id primitive.ObjectID, db *mongo.Database) (data []bson.M, err error) {
	var report []intermoni.Report
	collection := db.Collection("report")
	filter := bson.M{"mahasiswamagang._id": _id}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa context: %s", err)
	}
	for _, r := range report {
		mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(r.MahasiswaMagang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa magang: %s", err)
		}
		magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
		if err != nil {
			return data, fmt.Errorf("error GetAllReportByMahasiswa get magang: %s", err)
		}
		datareport := bson.M{
			"magang": magang,
			"_id": r.ID,
			"task": r.Task,
			"deskripsi": r.Deskripsi,
			"hasil": r.Hasil,
			"kehadiran": r.Kehadiran,
			"createdat": r.CreatedAt,
			"updatedat": r.UpdatedAt,
		}
		data = append(data, datareport)
	}
	return data, nil
}

func GetReportByID(_id primitive.ObjectID, db *mongo.Database) (data bson.M, err error) {
	var report intermoni.Report
	collection := db.Collection("report")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get mahasiswa magang: %s", err)
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get mahasiswa: %s", err)
	}
	magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get magang: %s", err)
	}
	data = bson.M{
		"mahasiswa": mahasiswa,
		"magang": magang,
		"_id": report.ID,
		"task": report.Task,
		"deskripsi": report.Deskripsi,
		"hasil": report.Hasil,
		"kehadiran": report.Kehadiran,
		"createdat": report.CreatedAt,
		"updatedat": report.UpdatedAt,
	}
	return data, nil
}

func GetReportByIDOlehPembimbing(_id primitive.ObjectID, db *mongo.Database) (data bson.M, err error) {
	var report intermoni.Report
	collection := db.Collection("report")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get mahasiswa magang: %s", err)
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get mahasiswa: %s", err)
	}
	magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
	if err != nil {
		return data, fmt.Errorf("error GetReportByID get magang: %s", err)
	}
	data = bson.M{
		"mahasiswa": mahasiswa,
		"magang": magang,
		"_id": report.ID,
		"task": report.Task,
		"deskripsi": report.Deskripsi,
		"kehadiran": report.Kehadiran,
		"createdat": report.CreatedAt,
		"updatedat": report.UpdatedAt,
	}
	return data, nil
}