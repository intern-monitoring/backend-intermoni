package report

import (
	"fmt"
	"time"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TambahReportByMahasiswa(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Report) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	if mahasiswa.ID != mahasiswa_magang.Mahasiswa.ID {
		return fmt.Errorf("kamu bukan mahasiswa magang ini")
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("kamu belum melakukan kontrak magang")
	}
	if insertedDoc.Penerima != mahasiswa_magang.Pembimbing.ID || insertedDoc.Penerima != mahasiswa_magang.Mentor.ID {
		return fmt.Errorf("kamu tidak dapat memberikan report selain kepada pembimbing dan mentor")
	}
	if insertedDoc.Judul == "" || insertedDoc.Isi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	insertedDoc.MahasiswaMagang = mahasiswa_magang
	insertedDoc.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	_, err = intermoni.InsertOneDoc(db, "report", insertedDoc)
	if err != nil {
		return err
	}
	return nil
}