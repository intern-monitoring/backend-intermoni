package email

import (
	"fmt"

	"github.com/badoux/checkmail"
	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mahasiswa & mitra
func UpdateEmailUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.User) error {
	dataUser, err := intermoni.GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	existsDoc, _ := intermoni.GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	user := bson.M{
		"email":    insertedDoc.Email,
		"password": dataUser.Password,
		"salt":     dataUser.Salt,
		"role":     dataUser.Role,
	}
	err = intermoni.UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}