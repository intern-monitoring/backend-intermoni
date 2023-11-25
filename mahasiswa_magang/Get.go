package mahasiswa_magang

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"github.com/intern-monitoring/backend-intermoni/magang"
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
		mahasiswa, err := intermoni.GetMahasiswaFromID(m.Mahasiswa.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mahasiswa: %s", err)
		}
		m.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromID(m.Magang.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get magang: %s", err)
		}
		m.Magang = magang
		pembimbing, _ := intermoni.GetPembimbingFromID(m.Pembimbing.ID, db)
		m.Pembimbing = pembimbing
		mentor, _ := intermoni.GetMentorFromID(m.Mentor.ID, db)
		m.Mentor = mentor
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
		mahasiswa, err := intermoni.GetMahasiswaFromID(m.Mahasiswa.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa get mahasiswa: %s", err)
		}
		m.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromID(m.Magang.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByMahasiswa get magang: %s", err)
		}
		m.Magang = magang
		pembimbing, _ := intermoni.GetPembimbingFromID(m.Pembimbing.ID, db)
		m.Pembimbing = pembimbing
		mentor, _ := intermoni.GetMentorFromID(m.Mentor.ID, db)
		m.Mentor = mentor
		mahasiswa_magang = append(mahasiswa_magang, m)
		mahasiswa_magang = mahasiswa_magang[1:]
	}
	return mahasiswa_magang, nil
}

// by mitra
func GetMahasiswaMagangByMitra(iduser primitive.ObjectID, db *mongo.Database) (mahasiswa_magang []intermoni.MahasiswaMagang, err error) {
	magang, err := magang.GetAllMagangByMitra(iduser, db)
	if err != nil {
		return mahasiswa_magang, err
	}
	for _, m := range magang {
		mahasiswa_mgn, err := GetAllMahasiswaMagangFromMagangByMitra(m.ID, iduser, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangByAdmin get mahasiswa: %s", err)
		}
		mahasiswa_magang = append(mahasiswa_magang, mahasiswa_mgn...)
	}
	return mahasiswa_magang, nil
}

func GetAllMahasiswaMagangFromMagangByMitra(idmagang, iduser primitive.ObjectID, db *mongo.Database) (mahasiswa_magang []intermoni.MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"magang._id": idmagang}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetAllMahasiswaMagangFromMagangByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetAllMahasiswaMagangFromMagangByMitra context: %s", err)
	}
	for _, m := range mahasiswa_magang {
		mahasiswa, err := intermoni.GetMahasiswaFromID(m.Mahasiswa.ID, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetAllMahasiswaMagangFromMagangByMitra get mahasiswa: %s", err)
		}
		m.Mahasiswa = mahasiswa
		magang, err := intermoni.GetMagangFromIDByMitra(m.Magang.ID, iduser, db)
		if err != nil {
			return mahasiswa_magang, fmt.Errorf("error GetAllMahasiswaMagangFromMagangByMitra get magang: %s", err)
		}
		m.Magang = magang
		mentor, _ := intermoni.GetMentorFromID(m.Mentor.ID, db)
		m.Mentor = mentor
		mahasiswa_magang = append(mahasiswa_magang, m)
		mahasiswa_magang = mahasiswa_magang[1:]
	}
	return mahasiswa_magang, nil
}

// get by id mahasiswa magang
func GetMahasiswaMagangFromID(idmahasiswamagang primitive.ObjectID, db *mongo.Database) (mahasiswa_magang intermoni.MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"_id": idmahasiswamagang}
	err = collection.FindOne(context.Background(), filter).Decode(&mahasiswa_magang)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID: %s", err)
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get mahasiswa: %s", err)
	}
	mahasiswa_magang.Mahasiswa = mahasiswa
	magang, err := intermoni.GetMagangFromID(mahasiswa_magang.Magang.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get magang: %s", err)
	}
	mahasiswa_magang.Magang = magang
	pembimbing, _ := intermoni.GetPembimbingFromID(mahasiswa_magang.Pembimbing.ID, db)
	mahasiswa_magang.Pembimbing = pembimbing
	mentor, _ := intermoni.GetMentorFromID(mahasiswa_magang.Mentor.ID, db)
	mahasiswa_magang.Mentor = mentor
	return mahasiswa_magang, nil
}