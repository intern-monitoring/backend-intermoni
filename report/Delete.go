package report

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteReportByMahasiswa(idreport, iduser primitive.ObjectID, db *mongo.Database) error {
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
	if mahasiswa_magang.Mahasiswa.ID != mahasiswa.ID {
		return fmt.Errorf("kamu bukan pemilik report ini")
	}
	err = intermoni.DeleteOneDoc(idreport, db, "report")
	if err != nil {
		return err
	}
	return nil
}