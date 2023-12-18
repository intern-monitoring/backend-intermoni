package mentor

import (
	"context"
	"fmt"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateMentor(idparam, iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
	mentor, err := intermoni.GetMentorFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if CheckMentor_MahasiswaMagang(mentor.ID, db) {
		return fmt.Errorf("kamu masih memiliki mahasiswa magang")
	}
	if mentor.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	namalengkap := r.FormValue("namalengkap")
	nik := r.FormValue("nik")

	if namalengkap == "" || nik == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	
	imageUrl, err := intermoni.SaveFileToGithub("Fatwaff", "fax.mp4@gmail.com", "bk-image", "user" ,r)
	if err != nil {
		return fmt.Errorf("error save file: %s", err)
	}

	data := bson.M{
		"namalengkap": namalengkap,
		"nik":         nik,
		"mitra": bson.M{
			"_id": mentor.Mitra.ID,
		},
		"imageurl": imageUrl,
		"akun": intermoni.User{
			ID: mentor.Akun.ID,
		},
	}
	err = intermoni.UpdateOneDoc(idparam, db, "mentor", data)
	if err != nil {
		return err
	}
	return nil
}

func CheckMentor_MahasiswaMagang(idmentor primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mentor._id": idmentor,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(idmentor, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idmentor primitive.ObjectID, db *mongo.Database) int64 {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mentor._id": idmentor,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	return count
}