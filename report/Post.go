package report

import (
	"context"
	"fmt"
	"net/http"
	"time"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TambahReportByMahasiswa(iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
	mahasiswa_magang, err := GetMahasiswaMagangByMahasiswa(iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("kamu belum lolos seleksi")
	}

	task := r.FormValue("task")
	deskripsi := r.FormValue("deskripsi")
	hasil := r.FormValue("hasil")

	if task == "" || deskripsi == "" || hasil == ""  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	imageUrl, err := intermoni.SaveFileToGithub("Fatwaff", "fax.mp4@gmail.com", "bk-image", "report" ,r)
	if err != nil {
		return fmt.Errorf("error save file: %s", err)
	}
	
	data := bson.M{
		"mahasiswamagang": bson.M{
			"_id": mahasiswa_magang.ID,
		},
		"task": task,
		"deskripsi": deskripsi,
		"hasil": hasil,
		"kehadiran": imageUrl,
		"createdat": primitive.NewDateTimeFromTime(time.Now().UTC()),
		"updatedat": primitive.NewDateTimeFromTime(time.Now().UTC()),
		"feedback": "",
		"nilaimentor": 0,
		"nilaipembimbing": 0,
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