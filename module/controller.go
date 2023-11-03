package module

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/intern-monitoring/backend-intermoni/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

// var MongoString string = os.Getenv("MONGOSTRING")

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// signup
func SignUpMahasiswa(db *mongo.Database, insertedDoc model.Mahasiswa) error {
	objectId := primitive.NewObjectID() 
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NIM == "" || insertedDoc.PerguruanTinggi == "" || insertedDoc.Prodi == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
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
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "mahasiswa",
	}
	mahasiswa := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nim": insertedDoc.NIM,
		"perguruantinggi": insertedDoc.PerguruanTinggi,
		"prodi": insertedDoc.Prodi,
		"akun": objectId,
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = InsertOneDoc(db, "mahasiswa", mahasiswa)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	return nil
}

func SignUpMitra(db *mongo.Database, insertedDoc model.Mitra) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaNarahubung == "" || insertedDoc.NoHpNarahubung == "" || insertedDoc.Nama == "" || insertedDoc.Kategori == "" || insertedDoc.SektorIndustri == "" || insertedDoc.Alamat == "" || insertedDoc.Website == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
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
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "mitra",
	}
	mitra := bson.M{
		"namanarahubung": insertedDoc.NamaNarahubung,
		"nohpnarahubung": insertedDoc.NoHpNarahubung,
		"namaresmi": insertedDoc.Nama,
		"kategori": insertedDoc.Kategori,
		"sektorindustri": insertedDoc.SektorIndustri,
		"alamat": insertedDoc.Alamat,
		"website": insertedDoc.Website,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	_, err = InsertOneDoc(db, "mitra", mitra)
	if err != nil {
		return err
	}
	return nil
}

// login
func LogIn(db *mongo.Database, insertedDoc model.User) (user model.User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("mohon untuk melengkapi data")
	} 
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("email tidak valid")
	} 
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return 
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

//user
func UpdateUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return err
	}
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if insertedDoc.Confirmpassword != insertedDoc.Password {
		return fmt.Errorf("konfirmasi password salah")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email": insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser(db *mongo.Database) (user []model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser mongo: %s", err)
	}
	err = cursor.All(context.Background(), &user)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser context: %s", err)
	}
	return user, nil
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// mahasiswa
// func UpdateMahasiswa(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Mahasiswa) error {
// 	mahasiswa, err := GetMahasiswaFromAkun(iduser, db)
// }

func GetAllMahasiswa(db *mongo.Database) (mahasiswa []model.Mahasiswa, err error) {
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

func GetMahasiswaFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetMahasiswaFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("mahasiswa tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// mitra
func GetMitraFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetMitraFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("mitra tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// magang
func InsertMagang(_id primitive.ObjectID, db *mongo.Database, insertedDoc model.Magang) error {
	if insertedDoc.Posisi == "" || insertedDoc.Lokasi == "" || insertedDoc.DeskripsiMagang == "" || insertedDoc.InfoTambahanMagang == "" || insertedDoc.Expired == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mitra, err := GetMitraFromAkun(_id, db)
	if err != nil {
		return err
	}
	magang := bson.M{
		"posisi": insertedDoc.Posisi,
		"mitra": model.Mitra {
			ID : mitra.ID,
			Akun: model.User{
				ID : _id,
			},
		},
		"lokasi": insertedDoc.Lokasi,
		"createdat": primitive.NewDateTimeFromTime(time.Now().UTC()),
		"deskripsimagang": insertedDoc.DeskripsiMagang,
		"infotambahanmagang": insertedDoc.InfoTambahanMagang,
		"expired": insertedDoc.Expired,
	}
	_, err = InsertOneDoc(db, "magang", magang)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMagang(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Magang) error {
	_, err := GetMagangFromIDByMitra(idparam, iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Posisi == "" || insertedDoc.Lokasi == "" || insertedDoc.DeskripsiMagang == "" || insertedDoc.InfoTambahanMagang == "" || insertedDoc.Expired == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	mitra, err := GetMitraFromAkun(iduser, db)
	if err != nil {
		return err
	}
	magang := bson.M{
		"posisi": insertedDoc.Posisi,
		"mitra": model.Mitra {
			ID : mitra.ID,
			Akun: model.User{
				ID : iduser,
			},
		},
		"lokasi": insertedDoc.Lokasi,
		"deskripsimagang": insertedDoc.DeskripsiMagang,
		"infotambahanmagang": insertedDoc.InfoTambahanMagang,
		"expired": insertedDoc.Expired,
	}
	err = UpdateOneDoc(idparam, db, "magang", magang)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMagang(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetMagangFromIDByMitra(idparam, iduser, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "magang")
	if err != nil {
		return err
	}
	return nil
}

func GetAllMagang(db *mongo.Database) (magang []model.Magang, err error) {
	collection := db.Collection("magang")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagang mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &magang)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagang context: %s", err)
	}
	for _, m := range magang {
		mitra, err := GetMitraFromID(m.Mitra.ID, db)
		if err != nil {
			fmt.Println(m.Mitra.ID)
			return magang, fmt.Errorf("error GetAllMagang get mitra: %s", err)
		}
		m.Mitra = mitra
		magang = append(magang, m)
		magang = magang[1:]
	}
	return magang, nil
}

func GetMagangFromMitra(_id primitive.ObjectID, db *mongo.Database) (magang []model.Magang, err error) {
	collection := db.Collection("magang")
	mitra, err := GetMitraFromAkun(_id, db)
	if err != nil {
		return magang, err
	}
	filter := bson.M{"mitra._id": mitra.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return magang, fmt.Errorf("error GetMagangByMitra mongo: %s", err)
	}
	err = cursor.All(context.Background(), &magang)
	if err != nil {
		return magang, fmt.Errorf("error GetMagangByMitra context: %s", err)
	}
	return magang, nil
}

func GetMagangFromID(_id primitive.ObjectID, db *mongo.Database) (magang model.Magang, err error) {
	collection := db.Collection("magang")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&magang)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return magang, fmt.Errorf("no data found for ID %s", _id)
		}
		return magang, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	mitra, err := GetMitraFromID(magang.Mitra.ID, db)
	if err != nil {
		return magang, fmt.Errorf("error GetAllMagang get mitra: %s", err)
	}
	magang.Mitra = mitra
	return magang, nil
}

func GetMagangFromIDByMitra(idparam, iduser primitive.ObjectID, db *mongo.Database) (magang model.Magang, err error) {
	collection := db.Collection("magang")
	filter := bson.M{"_id": idparam}
	err = collection.FindOne(context.TODO(), filter).Decode(&magang)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return magang, fmt.Errorf("no data found for ID %s", idparam)
		}
		return magang, fmt.Errorf("error retrieving data for ID %s: %s", idparam, err.Error())
	}
	mitra, err := GetMitraFromAkun(iduser, db)
	if err != nil {
		return magang, err
	}
	if magang.Mitra.ID != mitra.ID {
		return magang, fmt.Errorf("kamuh bukan pemilik magang ini")
	}
	return magang, nil
}