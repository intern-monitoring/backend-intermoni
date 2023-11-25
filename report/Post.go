package report

import (
	"context"
	"fmt"
	"time"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TambahReportByMahasiswa(iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Report) error {
	mahasiswa_magang, err := GetMahasiswaMagangByMahasiswa(iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("kamu belum lolos seleksi")
	}
	if insertedDoc.Judul == "" || insertedDoc.Isi == "" || insertedDoc.Penerima.ID == primitive.NilObjectID {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if insertedDoc.Penerima.ID != mahasiswa_magang.Pembimbing.ID && insertedDoc.Penerima.ID != mahasiswa_magang.Mentor.ID {
		return fmt.Errorf("kamu tidak dapat memberikan report selain kepada pembimbing dan mentor kamu")
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"judul":    insertedDoc.Judul,
		"isi":      insertedDoc.Isi,
		"penerima": bson.M{
			"_id": insertedDoc.Penerima.ID,
		},
		"createdat": primitive.NewDateTimeFromTime(time.Now().UTC()),
	}
	_, err = intermoni.InsertOneDoc(db, "report", data)
	if err != nil {
		return err
	}
	return nil
}

func GetMahasiswaMagangByMahasiswa(iduser primitive.ObjectID, db *mongo.Database) (mahasiswa_magang intermoni.MahasiswaMagang, err error) {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa get mahasiswa: %s", err)
	}
	filter := bson.M{"mahasiswa._id": mahasiswa.ID}
	err = db.Collection("mahasiswa_magang").FindOne(context.TODO(), filter).Decode(&mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa context: %s", err)
	}
	return mahasiswa_magang, nil
}