package mahasiswa_magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by admin
func GetMahasiswaMagangByAdmin(db *mongo.Database) (mahasiswa_magang []intermoni.MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin context: %s", err)
	}
	for _, m := range mahasiswa_magang {
		mahasiswa, err := intermoni.GetMahasiswaFromAkun(m.Mahasiswa.Akun.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mahasiswa: %s", err)
		}
		m.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromID(m.Magang.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get magang: %s", err)
		}
		mitra, err := intermoni.GetMitraFromAkun(magang.Mitra.Akun.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mitra: %s", err)
		}
		magang.Mitra = mitra
		m.Magang = magang
		mahasiswa_magang = append(mahasiswa_magang, m)
		mahasiswa_magang = mahasiswa_magang[1:]
	}
	return mahasiswa_magang, nil
}

// by mitra
func GetMahasiswaMagangByMitra(_id primitive.ObjectID, db *mongo.Database) (mahasiswa_magang []intermoni.MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	mitra, err := intermoni.GetMitraFromAkun(_id, db)
	if err != nil {
		return mahasiswa_magang, err
	}
	filter := bson.M{"magang.mitra._id": mitra.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMitra context: %s", err)
	}
	for _, m := range mahasiswa_magang {
		mahasiswa, err := intermoni.GetMahasiswaFromAkun(m.Mahasiswa.Akun.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mahasiswa: %s", err)
		}
		m.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromIDByMitra(m.Magang.ID, m.Magang.Mitra.Akun.ID,  db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get magang: %s", err)
		}
		m.Magang = magang
		mahasiswa_magang = append(mahasiswa_magang, m)
		mahasiswa_magang = mahasiswa_magang[1:]
	}
	return mahasiswa_magang, nil
}

// by mahasiswa
func GetMahasiswaMagangByMahasiswa(_id primitive.ObjectID, db *mongo.Database) (mahasiswa_magang []intermoni.MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(_id, db)
	if err != nil {
		return mahasiswa_magang, err
	}
	filter := bson.M{"mahasiswa._id": mahasiswa.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa context: %s", err)
	}
	for _, m := range mahasiswa_magang {
		magang, err := intermoni.GetMagangFromID(m.Magang.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get magang: %s", err)
		}
		mitra, err := intermoni.GetMitraFromAkun(magang.Mitra.Akun.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mitra: %s", err)
		}
		magang.Mitra = mitra
		m.Magang = magang
		mahasiswa_magang = append(mahasiswa_magang, m)
		mahasiswa_magang = mahasiswa_magang[1:]
	}
	return mahasiswa_magang, nil
}
