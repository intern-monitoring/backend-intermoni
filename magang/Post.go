package magang

import (
	"fmt"
	"time"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func TambahMagangOlehMitra(_id primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Magang) error {
	if insertedDoc.Posisi == "" || insertedDoc.Lokasi == "" || insertedDoc.DeskripsiMagang == "" || insertedDoc.InfoTambahanMagang == "" || insertedDoc.Expired == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mitra, err := intermoni.GetMitraFromAkun(_id, db)
	if err != nil {
		return err
	}
	magang := bson.M{
		"posisi": insertedDoc.Posisi,
		"mitra": bson.M{
			"_id": mitra.ID,
		},
		"lokasi":             insertedDoc.Lokasi,
		"createdat":          primitive.NewDateTimeFromTime(time.Now().UTC()),
		"deskripsimagang":    insertedDoc.DeskripsiMagang,
		"infotambahanmagang": insertedDoc.InfoTambahanMagang,
		"expired":            insertedDoc.Expired,
	}
	_, err = intermoni.InsertOneDoc(db, "magang", magang)
	if err != nil {
		return err
	}
	return nil
}