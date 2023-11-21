package magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func EditMagangOlehMitra(idmagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Magang) error {
	if CheckMagang_MahasiswaMagang(idmagang, db) {
		return fmt.Errorf("magang masih dalam proses seleksi")
	}
	magang, err := intermoni.GetMagangFromIDByMitra(idmagang, iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Posisi == "" || insertedDoc.Lokasi == "" || insertedDoc.DeskripsiMagang == "" || insertedDoc.InfoTambahanMagang == "" || insertedDoc.Expired == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	data := bson.M{
		"posisi": insertedDoc.Posisi,
		"mitra": bson.M{
			"_id": magang.Mitra.ID,
		},
		"lokasi":             insertedDoc.Lokasi,
		"deskripsimagang":    insertedDoc.DeskripsiMagang,
		"infotambahanmagang": insertedDoc.InfoTambahanMagang,
		"expired":            insertedDoc.Expired,
	}
	err = intermoni.UpdateOneDoc(idmagang, db, "magang", data)
	if err != nil {
		return err
	}
	return nil
}

func CheckMagang_MahasiswaMagang(idmagang primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"magang._id": idmagang,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(idmagang, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idmagang primitive.ObjectID, db *mongo.Database) int64 {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"magang._id": idmagang,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	return count
}