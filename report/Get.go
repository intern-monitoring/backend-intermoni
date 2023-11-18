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
func GetAllReportByMahasiswa(_id primitive.ObjectID, db *mongo.Database) (report intermoni.Report, err error) {
	collection := db.Collection("report")
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(_id, db)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa: %s", err)
	}
	filter := bson.M{"mahasiswa_magang.mahasiswa._id": mahasiswa.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMahasiswa context: %s", err)
	}
	return report, nil
}

// by mentor/pembimbing
func GetAllReportByPenerima(_id primitive.ObjectID, db *mongo.Database) (report []intermoni.Report, err error) {
	collection := db.Collection("report")
	filter := bson.M{"penerima": _id}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMitra context: %s", err)
	}
	return report, nil
}