package mahasiswa_magang

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mahasiswa
func BatalApply(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.SeleksiBerkas != 0 || mahasiswa_magang.SeleksiWewancara != 0 {
		return fmt.Errorf("kamu sudah dalam proses seleksi")
	}
	if mahasiswa_magang.Mahasiswa.ID != mahasiswa.ID {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	err = intermoni.DeleteOneDoc(idmahasiswamagang, db, "mahasiswa_magang")
	if err != nil {
		return err
	}
	return nil
}