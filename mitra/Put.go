package mitra

import (
	"context"
	"fmt"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"github.com/intern-monitoring/backend-intermoni/mahasiswa_magang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var imageUrl string

// by mitra
func UpdateMitra(idparam, iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
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

	namanarahubung := r.FormValue("namanarahubung")
	nohpnarahubung := r.FormValue("nohpnarahubung")
	nama := r.FormValue("nama")
	kategori := r.FormValue("kategori")
	sektorindustri := r.FormValue("sektorindustri")
	tentang := r.FormValue("tentang")
	alamat := r.FormValue("alamat")
	website := r.FormValue("website")
	img := r.FormValue("file")


	if namanarahubung == "" || nohpnarahubung == "" || nama == "" || kategori == "" || sektorindustri == "" || alamat == "" || website == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	if img != "" {
		imageUrl = img
	} else {
		imageUrl, err = intermoni.SaveFileToGithub("Fatwaff", "fax.mp4@gmail.com", "bk-image", "user" ,r)
		if err != nil {
			return fmt.Errorf("error save file: %s", err)
		}
	}

	mtr := bson.M{
		"namanarahubung": namanarahubung,
		"nohpnarahubung": nohpnarahubung,
		"nama":           nama,
		"kategori":       kategori,
		"sektorindustri": sektorindustri,
		"tentang":        tentang,
		"alamat":         alamat,
		"website":        website,
		"mou":            0,
		"imageurl":       imageUrl,
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