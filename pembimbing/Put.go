package pembimbing

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdatePembimbing(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Pembimbing) error {
	pembimbing, err := intermoni.GetPembimbingFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if pembimbing.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.NIK == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	data := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"nik":         insertedDoc.NIK,
		"akun": intermoni.User{
			ID: pembimbing.Akun.ID,
		},
		"prodi": insertedDoc.Prodi,
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