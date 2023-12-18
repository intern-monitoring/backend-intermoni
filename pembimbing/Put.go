package pembimbing

import (
	"context"
	"fmt"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdatePembimbing(idparam, iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
	pembimbing, err := intermoni.GetPembimbingFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if pembimbing.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	namalengkap := r.FormValue("namalengkap")
	nik := r.FormValue("nik")
	prodi := r.FormValue("prodi")

	if namalengkap == "" || nik == "" || prodi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	imageUrl, err := intermoni.SaveFileToGithub("Fatwaff", "fax.mp4@gmail.com", "bk-image", "user" ,r)
	if err != nil {
		return fmt.Errorf("error save file: %s", err)
	}

	data := bson.M{
		"namalengkap": namalengkap,
		"nik":         nik,
		"prodi": prodi,
		"imageurl": imageUrl,
		"akun": intermoni.User{
			ID: pembimbing.Akun.ID,
		},
	}
	err = intermoni.UpdateOneDoc(idparam, db, "pembimbing", data)
	if err != nil {
		return err
	}
	return nil
}

func CheckMentor_MahasiswaMagang(idpembimbing primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"pembimbing._id": idpembimbing,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(idpembimbing, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idpembimbing primitive.ObjectID, db *mongo.Database) int64 {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"pembimbing._id": idpembimbing,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	return count
}