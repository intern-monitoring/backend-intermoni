package signup_mahasiswa

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/badoux/checkmail"
	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func SignUpMahasiswa(db *mongo.Database, insertedDoc intermoni.Mahasiswa) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NIM == "" || insertedDoc.PerguruanTinggi == "" || insertedDoc.Prodi == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := intermoni.GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if insertedDoc.Akun.Confirmpassword != insertedDoc.Akun.Password {
		return fmt.Errorf("konfirmasi password salah")
	}
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id":      objectId,
		"email":    insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     "mahasiswa",
	}
	mahasiswa := bson.M{
		"namalengkap":     insertedDoc.NamaLengkap,
		"tanggallahir":    insertedDoc.TanggalLahir,
		"jeniskelamin":    insertedDoc.JenisKelamin,
		"nim":             insertedDoc.NIM,
		"perguruantinggi": insertedDoc.PerguruanTinggi,
		"prodi":           insertedDoc.Prodi,
		"seleksikampus":   0,
		"akun": intermoni.User{
			ID: objectId,
		},
	}
	_, err = intermoni.InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = intermoni.InsertOneDoc(db, "mahasiswa", mahasiswa)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	return nil
}