package intermoni

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo
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

// get user
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc User, err error) {
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

func GetUserFromEmail(email string, db *mongo.Database) (doc User, err error) {
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

func GetUserFromPhone(phone string, db *mongo.Database) (doc User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"phone": phone}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("phone tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// get mahasiswa
func GetMahasiswaFromID(_id primitive.ObjectID, db *mongo.Database) (doc Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

func GetMahasiswaFromAkun(id_akun primitive.ObjectID, db *mongo.Database) (doc Mahasiswa, err error) {
	collection := db.Collection("mahasiswa")
	filter := bson.M{"akun._id": id_akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("mahasiswa tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

// get mitra
func GetMitraFromID(_id primitive.ObjectID, db *mongo.Database) (doc Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

func GetMitraFromAkun(id_akun primitive.ObjectID, db *mongo.Database) (doc Mitra, err error) {
	collection := db.Collection("mitra")
	filter := bson.M{"akun._id": id_akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("mitra tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

// get pembimbing
func GetPembimbingFromID(_id primitive.ObjectID, db *mongo.Database) (doc Pembimbing, err error) {
	collection := db.Collection("pembimbing")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

func GetPembimbingFromAkun(id_akun primitive.ObjectID, db *mongo.Database) (doc Pembimbing, err error) {
	collection := db.Collection("pembimbing")
	filter := bson.M{"akun._id": id_akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("pembimbing tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

// get mentor
func GetMentorFromID(_id primitive.ObjectID, db *mongo.Database) (doc Mentor, err error) {
	collection := db.Collection("mentor")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

func GetMentorFromAkun(id_akun primitive.ObjectID, db *mongo.Database) (doc Mentor, err error) {
	collection := db.Collection("mentor")
	filter := bson.M{"akun._id": id_akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("mentor tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	user, err := GetUserFromID(doc.Akun.ID, db)
	if err != nil {
		return doc, fmt.Errorf("kesalahan server")
	}
	akun := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	doc.Akun = akun
	return doc, nil
}

// get magang
func GetMagangFromID(_id primitive.ObjectID, db *mongo.Database) (magang Magang, err error) {
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
		return magang, fmt.Errorf("error GetMagang get mitra: %s", err)
	}
	magang.Mitra = mitra
	return magang, nil
}

func GetMagangFromIDByMitra(idparam, iduser primitive.ObjectID, db *mongo.Database) (magang Magang, err error) {
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

// get mahasiswa_magang
func GetDetailMahasiswaMagangFromID(_id primitive.ObjectID, db *mongo.Database) (mahasiswa_magang MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&mahasiswa_magang)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return mahasiswa_magang, fmt.Errorf("no data found for ID %s", _id)
		}
		return mahasiswa_magang, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	mahasiswa, err := GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get mahasiswa: %s", err)
	}
	mahasiswa_magang.Mahasiswa = mahasiswa
	magang, err := GetMagangFromID(mahasiswa_magang.Magang.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get magang: %s", err)
	}
	mahasiswa_magang.Magang = magang
	pembimbing, err := GetPembimbingFromID(mahasiswa_magang.Pembimbing.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get pembimbing: %s", err)
	}
	mahasiswa_magang.Pembimbing = pembimbing
	mentor, err := GetMentorFromID(mahasiswa_magang.Mentor.ID, db)
	if err != nil {
		return mahasiswa_magang, fmt.Errorf("error GetMahasiswaMagangFromID get mentor: %s", err)
	}
	mahasiswa_magang.Mentor = mentor
	return mahasiswa_magang, nil
}

func GetMahasiswaMagangFromID(_id primitive.ObjectID, db *mongo.Database) (mahasiswa_magang MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&mahasiswa_magang)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return mahasiswa_magang, fmt.Errorf("no data found for ID %s", _id)
		}
		return mahasiswa_magang, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return mahasiswa_magang, nil
}

func GetMahasiswaMagangByMahasiswa(idmahasiswa primitive.ObjectID, db *mongo.Database) (mahasiswa_magang MahasiswaMagang, err error) {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{"mahasiswa._id": idmahasiswa}
	err = collection.FindOne(context.Background(), filter).Decode(&mahasiswa_magang)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return mahasiswa_magang, fmt.Errorf("mahasiswa magang tidak ditemukan")
		}
		return mahasiswa_magang, fmt.Errorf("terjadi kesalahan")
	}
	return mahasiswa_magang, nil
}

// get report
func GetReportFromID(_id primitive.ObjectID, db *mongo.Database) (report Report, err error) {
	filter := bson.M{"_id": _id}
	err = db.Collection("report").FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return report, fmt.Errorf("no data found for ID %s", _id)
		}
		return report, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return report, nil
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (Payload, error) {
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// get id
func GetID(r *http.Request) string {
    return r.URL.Query().Get("id")
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// save file to github
func SaveFileToGithub(usernameGhp, emailGhp, repoGhp, path string, r *http.Request) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	// Generate a random filename
	randomFileName, err := generateRandomFileName(handler.Filename)
	if err != nil {
		return "", fmt.Errorf("error 2: %s", err)
	}

	// Read the content of the file into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error 5: %s", err)
	}

	access_token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if access_token == "" {
		return "", fmt.Errorf("error access token: %s", err)
	}

	// Initialize GitHub client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access_token},
	)
	tc := oauth2.NewClient(r.Context(), ts)
	client := github.NewClient(tc)

	// Create a new repository file
	_, _, err = client.Repositories.CreateFile(r.Context(), usernameGhp, repoGhp, path+"/"+randomFileName, &github.RepositoryContentFileOptions{
		Message:   github.String("Add new file"),
		Content:   fileContent,
		Committer: &github.CommitAuthor{Name: github.String(usernameGhp), Email: github.String(emailGhp)},
	})
	if err != nil {
		return "", fmt.Errorf("error 6: %s", err)
	}

	imageUrl := "https://raw.githubusercontent.com/" + usernameGhp + "/" + repoGhp + "/main/"+path+"/" + randomFileName

	return imageUrl, nil
}

func generateRandomFileName(originalFilename string) (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomFileName := fmt.Sprintf("%x%s", randomBytes, filepath.Ext(originalFilename))
	return randomFileName, nil
}

func SendWhatsAppConfirmation(iduser primitive.ObjectID, db *mongo.Database, message string) error {
	url := "https://api.wa.my.id/api/send/message/text"

	user, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}

	// Data yang akan dikirimkan dalam format JSON
	jsonStr := []byte(`{
        "to": "` + user.Phone + `",
        "isgroup": false,
        "messages": "` + message + `"
    }`)

	// Membuat permintaan HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	// Menambahkan header ke permintaan
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", os.Getenv("TOKENWEBHOOK"))
	// req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDI0LTAyLTE5VDIxOjA3OjM2WiIsImlhdCI6IjIwMjQtMDEtMjBUMjE6MDc6MzZaIiwiaWQiOiI2MjgyMzE3MTUwNjgxIiwibmJmIjoiMjAyNC0wMS0yMFQyMTowNzozNloiff1YQuHHPwSzGpisAMb9rTLP58-jCqtByzePJACBLghprkq2HXtTSbVTShc49m3GIVkU42VSl8uSGme8c4vXnQc")
	req.Header.Set("Content-Type", "application/json")

	// Melakukan permintaan HTTP POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Menampilkan respons dari server
	fmt.Println("Response Status:", resp.Status)

	return nil
}