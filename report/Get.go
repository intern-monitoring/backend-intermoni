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
func GetAllReportByMahasiswa(_id primitive.ObjectID, db *mongo.Database) (data bson.M, err error) {
	var report []intermoni.Report
	var penerima bson.M
	collection := db.Collection("report")
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(_id, db)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa: %s", err)
	}
	filter := bson.M{"mahasiswa_magang.mahasiswa._id": mahasiswa.ID}
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
		pembimbing, _ := intermoni.GetPembimbingFromID(r.Penerima.ID, db)
		mentor, _ := intermoni.GetMentorFromID(r.Penerima.ID, db)
		if pembimbing.ID != primitive.NilObjectID {
			penerima = bson.M{
				"_id": pembimbing.ID,
				"nama": pembimbing.NamaLengkap,
				"nik": pembimbing.NIK,
				"prodi": pembimbing.Prodi,
			}
		} else if mentor.ID != primitive.NilObjectID {
			penerima = bson.M{
				"_id": mentor.ID,
				"nama": mentor.NamaLengkap,
				"nik": mentor.NIK,
			}
		} else {
			penerima = bson.M{
				"_id": primitive.NilObjectID,
			}
		}
		report = append(report, r)
		report = report[1:]
		data = bson.M{
			"magang": magang,
			"report": report,
			"penerima": penerima,
		}
	}
	return data, nil
}

// by mentor/pembimbing
func GetAllReportByPenerima(_id primitive.ObjectID, db *mongo.Database) (report []intermoni.Report, err error) {
	collection := db.Collection("report")
	filter := bson.M{"penerima._id": _id}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &report)
	if err != nil {
		return report, fmt.Errorf("error GetAllReportByMitra context: %s", err)
	}
	for _, r := range report {
		mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(r.MahasiswaMagang.ID, db)
		if err != nil {
			return report, fmt.Errorf("error GetAllReportByMitra get mahasiswa magang: %s", err)
		}
		mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
		if err != nil {
			return report, fmt.Errorf("error GetAllReportByMitra get mahasiswa: %s", err)
		}
		mahasiswa_magang.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
		if err != nil {
			return report, fmt.Errorf("error GetAllReportByMitra get magang: %s", err)
		}
		mahasiswa_magang.Magang = magang
		report = append(report, r)
		report = report[1:]
	}
	return report, nil
}

func GetReportByID(_id primitive.ObjectID, db *mongo.Database) (report intermoni.Report, err error) {
	collection := db.Collection("report")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		return report, fmt.Errorf("error GetReportByID: %s", err)
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return report, fmt.Errorf("error GetReportByID get mahasiswa magang: %s", err)
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return report, fmt.Errorf("error GetReportByID get mahasiswa: %s", err)
	}
	mahasiswa_magang.Mahasiswa = mahasiswa
	magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
	if err != nil {
		return report, fmt.Errorf("error GetReportByID get magang: %s", err)
	}
	mahasiswa_magang.Magang = magang
	penerima, err := intermoni.GetUserFromID(report.Penerima.ID, db)
	if err != nil {
		return report, fmt.Errorf("error GetReportByID get penerima: %s", err)
	}
	akun := intermoni.User{
		ID:    penerima.ID,
		Email: penerima.Email,
	}
	report.Penerima = akun
	report.MahasiswaMagang = mahasiswa_magang
	return report, nil
}