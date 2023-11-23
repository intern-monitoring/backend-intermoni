package mahasiswa_magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mahasiswa
func ApplyMagang(idmagang, iduser primitive.ObjectID, db *mongo.Database) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa.SeleksiKampus != 1 {
		return fmt.Errorf("kamu belum lolos seleksi kampus")
	}
	magang, err := intermoni.GetMagangFromID(idmagang, db)
	if err != nil {
		return err
	}
	if CheckMahasiswaMagang(mahasiswa.ID, magang.ID, db) {
		return fmt.Errorf("kamu sudah apply magang ini")
	}
	mahasiswa_magang := bson.M{
		"mahasiswa": bson.M{
			"_id": mahasiswa.ID,
		},
		"magang": bson.M{
			"_id": magang.ID,
		},
		"pembimbing": bson.M{
			"_id": primitive.NilObjectID,
		},
		"mentor": bson.M{
			"_id": primitive.NilObjectID,
		},
		"seleksiberkas":    0,
		"seleksiwewancara": 0,
		"status": 0,
	}
	_, err = intermoni.InsertOneDoc(db, "mahasiswa_magang", mahasiswa_magang)
	if err != nil {
		return err
	}
	return nil
}

func CheckMahasiswaMagang(idmahasiswa, idmagang primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"mahasiswa._id": idmahasiswa, "magang._id": idmagang}
	err := collection.FindOne(context.Background(), filter).Decode(&intermoni.MahasiswaMagang{})
	return err == nil
}
