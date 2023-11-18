package seleksi_berkas

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func SeleksiWewancaraMahasiswaMagangByMitra(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	_, err = intermoni.GetMagangFromIDByMitra(mahasiswa_magang.Magang.ID, iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.SeleksiWewancara != 0 {
		return fmt.Errorf("mahasiswa sudah diseleksi")
	}
	if insertedDoc.SeleksiWewancara != 1 && insertedDoc.SeleksiWewancara != 2 {
		return fmt.Errorf("kesalahan server")
	}
	if mahasiswa_magang.SeleksiBerkas != 1 {
		return fmt.Errorf("belum lolos seleksi berkas")
	}
	if insertedDoc.SeleksiWewancara == 2 {
		mahasiswa_magang.Status = 2
	}
	mahasiswa_magang.SeleksiWewancara = insertedDoc.SeleksiWewancara
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", mahasiswa_magang)
	if err != nil {
		return err
	}
	return nil
}