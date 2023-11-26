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
func GetAllReportByMahasiswa(_id primitive.ObjectID, db *mongo.Database) (data []bson.M, err error) {
	var report []intermoni.Report
	var penerima bson.M
	collection := db.Collection("report")
	mahasiswa_magang, err := GetMahasiswaMagangByMahasiswa(_id, db)
	if err != nil {
		return data, fmt.Errorf("error GetAllReportByMahasiswa get mahasiswa: %s", err)
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
		datareport := bson.M{
			"magang": magang,
			"_id": r.ID,
			"judul": r.Judul,
			"isi": r.Isi,
			"penerima": penerima,
			"createdat": r.CreatedAt,
		}
		data = append(data, datareport)
	}
	return data, nil
}

// by mentor/pembimbing
func GetAllReportByPenerima(_id primitive.ObjectID, db *mongo.Database) (report []intermoni.Report, err error) {
	var penerimaid primitive.ObjectID
	pembimbing, _ := intermoni.GetPembimbingFromAkun(_id, db)
	mentor, _ := intermoni.GetMentorFromAkun(_id, db)
	if pembimbing.ID != primitive.NilObjectID {
		penerimaid = pembimbing.ID
	} else if mentor.ID != primitive.NilObjectID {
		penerimaid = mentor.ID
	} else {
		penerimaid = primitive.NilObjectID
	}
	collection := db.Collection("report")
	filter := bson.M{"penerima._id": penerimaid}
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
		r.MahasiswaMagang = mahasiswa_magang
		report = append(report, r)
		report = report[1:]
	}
	return report, nil
}

func GetReportByID(_id primitive.ObjectID, db *mongo.Database) (data bson.M, err error) {
	var report intermoni.Report
	var penerima bson.M
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
	pembimbing, _ := intermoni.GetPembimbingFromID(report.Penerima.ID, db)
	mentor, _ := intermoni.GetMentorFromID(report.Penerima.ID, db)
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
	data = bson.M{
		"mahasiswa": mahasiswa,
		"magang": magang,
		"_id": report.ID,
		"judul": report.Judul,
		"isi": report.Isi,
		"penerima": penerima,
		"createdat": report.CreatedAt,
	}
	return data, nil
}