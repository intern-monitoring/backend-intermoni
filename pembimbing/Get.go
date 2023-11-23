package pembimbing

import (
	"context"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// by admin
func GetAllPembimbingByAdmin(db *mongo.Database) (pembimbing []intermoni.Pembimbing, err error) {
	filter := bson.M{}
	cursor, err := db.Collection("pembimbing").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &pembimbing); err != nil {
		return nil, err
	}
	for _, p := range pembimbing {
		user, err := intermoni.GetUserFromID(p.Akun.ID, db)
		if err != nil {
			return nil, err
		}
		akun := intermoni.User{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}
		p.Akun = akun
		pembimbing = append(pembimbing, p)
		pembimbing = pembimbing[1:]
	}
	return pembimbing, nil
}