package mahasiswa

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mahasiswa
func UpdateMahasiswa(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mahasiswa) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if CheckMahasiswa_MahasiswaMagang(mahasiswa.ID, db) {
		return fmt.Errorf("kamu masih dalam proses seleksi/magang")
	}
	if mahasiswa.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NIM == "" || insertedDoc.PerguruanTinggi == "" || insertedDoc.Prodi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mhs := bson.M{
		"namalengkap":     insertedDoc.NamaLengkap,
		"tanggallahir":    insertedDoc.TanggalLahir,
		"jeniskelamin":    insertedDoc.JenisKelamin,
		"nim":             insertedDoc.NIM,
		"perguruantinggi": insertedDoc.PerguruanTinggi,
		"prodi":           insertedDoc.Prodi,
		"seleksikampus":   0,
		"akun": intermoni.User{
			ID: mahasiswa.Akun.ID,
		},
	}
	err = intermoni.UpdateOneDoc(idparam, db, "mahasiswa", mhs)
	if err != nil {
		return err
	}
	return nil
}

// by admin
func SeleksiMahasiswaByAdmin(idparam primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mahasiswa) error {
	mahasiswa, err := intermoni.GetMahasiswaFromID(idparam, db)
	if err != nil {
		return err
	}
	if CheckMahasiswa_MahasiswaMagang(mahasiswa.ID, db) {
		return fmt.Errorf("mahasiswa masih dalam proses seleksi")
	}
	if insertedDoc.SeleksiKampus != 1 && insertedDoc.SeleksiKampus != 2 {
		return fmt.Errorf("kesalahan server")
	}
	mahasiswa.SeleksiKampus = insertedDoc.SeleksiKampus
	err = intermoni.UpdateOneDoc(idparam, db, "mahasiswa", mahasiswa)
	if err != nil {
		return err
	}
	return nil
}

func CheckMahasiswa_MahasiswaMagang(idmahasiswa primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mahasiswa._id": idmahasiswa,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(idmahasiswa, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idmahasiswa primitive.ObjectID, db *mongo.Database) int64 {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mahasiswa._id": idmahasiswa,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	return count
}