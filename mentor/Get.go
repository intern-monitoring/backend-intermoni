package mentor

import (
	"context"
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func GetAllMentorByMitra(iduser primitive.ObjectID, db *mongo.Database) (mentor []intermoni.Mentor, err error) {
	mitra, err := intermoni.GetMitraFromAkun(iduser, db)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"mitra._id": mitra.ID}
	cursor, err := db.Collection("mentor").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &mentor); err != nil {
		return nil, err
	}
	for _, m := range mentor {
		user, err := intermoni.GetUserFromID(m.Akun.ID, db)
		if err != nil {
			return nil, fmt.Errorf("error GetAllMentorByMitra get user: %s", err)
		}
		akun := intermoni.User{
			ID: user.ID,
			Email: user.Email,
			Role: user.Role,
		}
		m.Akun = akun
		mentor = append(mentor, m)
		mentor = mentor[1:]
	}
	return mentor, nil
}

// by admin
func GetAllMentorByAdmin(db *mongo.Database) (mentor []intermoni.Mentor, err error) {
	filter := bson.M{}
	cursor, err := db.Collection("mentor").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &mentor); err != nil {
		return nil, err
	}
	for _, m := range mentor {
		user, err := intermoni.GetUserFromID(m.Akun.ID, db)
		if err != nil {
			return nil, fmt.Errorf("error GetAllMentorByAdmin get user: %s", err)
		}
		akun := intermoni.User{
			ID: user.ID,
			Email: user.Email,
			Role: user.Role,
		}
		m.Akun = akun
		mentor = append(mentor, m)
		mentor = mentor[1:]
	}
	return mentor, nil
}

func GetMentorFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (mentor intermoni.Mentor, err error) {
	filter := bson.M{"_id": idparam}
	err = db.Collection("mentor").FindOne(context.Background(), filter).Decode(&mentor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return mentor, fmt.Errorf("mentor tidak ditemukan")
		}
		return mentor, err
	}
	user, err := intermoni.GetUserFromID(mentor.Akun.ID, db)
	if err != nil {
		return mentor, fmt.Errorf("error GetMentorFromID mongo: %s", err)
	}
	mitra, err := intermoni.GetMitraFromID(mentor.Mitra.ID, db)
	if err != nil {
		return mentor, fmt.Errorf("error GetMentorFromID mongo: %s", err)
	}
	mentor.Mitra = mitra
	akun := intermoni.User{
		ID: user.ID,
		Email: user.Email,
		Role: user.Role,
	}
	mentor.Akun = akun
	return mentor, nil
}