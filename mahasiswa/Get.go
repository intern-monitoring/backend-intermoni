package mahasiswa

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by admin
func GetMahasiswaFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (mahasiswa intermoni.Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&mahasiswa)
	if err != nil {
		return mahasiswa, fmt.Errorf("error GetMahasiswaFromID mongo: %s", err)
	}
	user, err := intermoni.GetUserFromID(mahasiswa.Akun.ID, db)
	if err != nil {
		return mahasiswa, fmt.Errorf("error GetMahasiswaFromID mongo: %s", err)
	}
	akun := intermoni.User{
		ID: user.ID,
		Email: user.Email,
		Role: user.Role,
	}
	mahasiswa.Akun = akun
	return mahasiswa, nil
}

func GetAllMahasiswaByAdmin(db *mongo.Database) (mahasiswa []intermoni.Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return mahasiswa, fmt.Errorf("error GetAllMahasiswa mongo: %s", err)
	}
	err = cursor.All(context.Background(), &mahasiswa)
	if err != nil {
		return mahasiswa, fmt.Errorf("error GetAllMahasiswa context: %s", err)
	}
	return mahasiswa, nil
}