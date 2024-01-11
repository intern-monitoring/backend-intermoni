package report

import (
	"fmt"
	"net/http"
	"time"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var imageUrl string

func UpdateReportByMahasiswa(idreport, iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	report, err := intermoni.GetReportFromID(idreport, db)
	if err != nil {
		return err
	}
	if report.Feedback != "" || report.NilaiMentor != 0 || report.NilaiPembimbing != 0 {
		return fmt.Errorf("report ini sudah di feedback")
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Mahasiswa.ID != mahasiswa.ID {
		return fmt.Errorf("kamu bukan pemilik report ini")
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("kamu belum lolos seleksi")
	}

	task := r.FormValue("task")
	deskrisi := r.FormValue("deskripsi")
	hasil := r.FormValue("hasil")
	file := r.FormValue("file")

	if task == "" || deskrisi == "" || hasil == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	if file != "" {
		imageUrl = file
	} else {
		imageUrl, err = intermoni.SaveFileToGithub("Fatwaff", "fax.mp4@gmail.com", "bk-image", "report" ,r)
		if err != nil {
			return fmt.Errorf("error save file: %s", err)
		}
	}

	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"task": task,
		"deskripsi":   deskrisi,
		"hasil":       hasil,
		"kehadiran":   imageUrl,
		"createdat": report.CreatedAt,
		"updatedat": primitive.NewDateTimeFromTime(time.Now().UTC()),
		"feedback": "",
		"nilaimentor": 0,
		"nilaipembimbing": 0,
	}
	err = intermoni.UpdateOneDoc(idreport, db, "report", data)
	if err != nil {
		return err
	}
	return nil
}

func TambahFeedbackNilaiByMentor(idreport, iduser primitive.ObjectID, db *mongo.Database, updatedDoc intermoni.Report) error {
	if updatedDoc.Feedback == "" || updatedDoc.NilaiMentor == 0 {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mentor, err := intermoni.GetMentorFromAkun(iduser, db)
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
	if mahasiswa_magang.Mentor.ID != mentor.ID {
		return fmt.Errorf("kamu bukan mentor dari report ini")
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"task": report.Task,
		"deskripsi":   report.Deskripsi,
		"hasil":       report.Hasil,
		"kehadiran":   report.Kehadiran,
		"createdat": report.CreatedAt,
		"updatedat": report.UpdatedAt,
		"feedback": updatedDoc.Feedback,
		"nilaimetor": updatedDoc.NilaiMentor,
		"nilaipembimbing": report.NilaiPembimbing,
	}
	err = intermoni.UpdateOneDoc(idreport, db, "report", data)
	if err != nil {
		return err
	}
	return nil
}

func TambahNilaiByPembimbing(idreport, iduser primitive.ObjectID, db *mongo.Database, updatedDoc intermoni.Report) error {
	if updatedDoc.NilaiPembimbing == 0 {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	pembimbing, err := intermoni.GetPembimbingFromAkun(iduser, db)
	if err != nil {
		return err
	}
	report, err := intermoni.GetReportFromID(idreport, db)
	if err != nil {
		return err
	}
	if report.Feedback == "" || report.NilaiMentor == 0 {
		return fmt.Errorf("mentor belum memberi feedback")
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(report.MahasiswaMagang.ID, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Pembimbing.ID != pembimbing.ID {
		return fmt.Errorf("kamu bukan pembimbing dari report ini")
	}
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"task": report.Task,
		"deskripsi":   report.Deskripsi,
		"hasil":       report.Hasil,
		"kehadiran":   report.Kehadiran,
		"createdat": report.CreatedAt,
		"updatedat": report.UpdatedAt,
		"feedback": report.Feedback,
		"nilaimetor": report.NilaiMentor,
		"nilaipembimbing": updatedDoc.NilaiPembimbing,
	}
	err = intermoni.UpdateOneDoc(idreport, db, "report", data)
	if err != nil {
		return err
	}
	return nil
}