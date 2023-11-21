package mitra

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"github.com/intern-monitoring/backend-intermoni/mahasiswa_magang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func UpdateMitra(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mitra) error {
	mitra, err := intermoni.GetMitraFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if CheckMitra_MahasiswaMagang(iduser, db) {
		return fmt.Errorf("kamu masih dalam proses seleksi/magang")
	}
	if mitra.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	if insertedDoc.NamaNarahubung == "" || insertedDoc.NoHpNarahubung == "" || insertedDoc.Nama == "" || insertedDoc.Kategori == "" || insertedDoc.SektorIndustri == "" || insertedDoc.Alamat == "" || insertedDoc.Website == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mtr := bson.M{
		"namanarahubung": insertedDoc.NamaNarahubung,
		"nohpnarahubung": insertedDoc.NoHpNarahubung,
		"nama":           insertedDoc.Nama,
		"kategori":       insertedDoc.Kategori,
		"sektorindustri": insertedDoc.SektorIndustri,
		"tentang":        insertedDoc.Tentang,
		"alamat":         insertedDoc.Alamat,
		"website":        insertedDoc.Website,
		"mou":            0,
		"akun": intermoni.User{
			ID: mitra.Akun.ID,
		},
	}
	err = intermoni.UpdateOneDoc(idparam, db, "mitra", mtr)
	if err != nil {
		return err
	}
	return nil
}

// by admin
func ConfirmMouMitraByAdmin(idparam primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mitra) error {
	mitra, err := intermoni.GetMitraFromID(idparam, db)
	if err != nil {
		return err
	}
	if CheckMitra_MahasiswaMagang(mitra.ID, db) {
		return fmt.Errorf("mitra masih dalam proses seleksi")
	}
	if insertedDoc.MoU != 1 && insertedDoc.MoU != 2 {
		return fmt.Errorf("kesalahan server")
	}
	mitra.MoU = insertedDoc.MoU
	err = intermoni.UpdateOneDoc(idparam, db, "mitra", mitra)
	if err != nil {
		return err
	}
	return nil
}

func CheckMitra_MahasiswaMagang(iduser primitive.ObjectID, db *mongo.Database) bool {
	mitra, _ := intermoni.GetMitraFromAkun(iduser, db)
	mahasiswa_magang, _ := mahasiswa_magang.GetMahasiswaMagangByMitra(iduser, db)
	count := len(mahasiswa_magang)
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(mitra.ID, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idmitra primitive.ObjectID, db *mongo.Database) int {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mitra._id": idmitra,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	countInt := int(count)
	return countInt
}