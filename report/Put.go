package report

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateReportByMahasiswa(idreport, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Report) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	report, err := intermoni.GetReportFromID(idreport, db)
	if err != nil {
		return err
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return err
	}
	pembimbing, err := intermoni.GetPembimbingFromAkun(insertedDoc.Penerima.ID, db)
	if err != nil {
		return err
	}
	mentor, err := intermoni.GetMentorFromAkun(insertedDoc.Penerima.ID, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Mahasiswa.ID != mahasiswa.ID {
		return fmt.Errorf("kamu bukan pemilik report ini")
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("kamu belum lolos seleksi")
	}
	if pembimbing.ID != mahasiswa_magang.Pembimbing.ID && mentor.ID != mahasiswa_magang.Mentor.ID {
		return fmt.Errorf("kamu tidak dapat memberikan report selain kepada pembimbing dan mentor kamu")
	}
	if insertedDoc.Judul == "" || insertedDoc.Isi == "" || insertedDoc.Penerima.ID == primitive.NilObjectID {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"judul": insertedDoc.Judul,
		"isi":   insertedDoc.Isi,
		"penerima": bson.M{
			"_id": insertedDoc.Penerima.ID,
		},
		"createdat": report.CreatedAt,
	}
	err = intermoni.UpdateOneDoc(idreport, db, "report", data)
	if err != nil {
		return err
	}
	return nil
}