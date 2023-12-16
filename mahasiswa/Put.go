package mahasiswa

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// const (
// 	githubUser     = "Fatwaff"
// 	repoName       = "bk-image"
// 	accessToken    = "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"
// 	uploadsDirPath = "user"
// )

// by mahasiswa
func UpdateMahasiswa(idparam, iduser primitive.ObjectID, db *mongo.Database, r *http.Request) error {
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if CheckMahasiswa_MahasiswaMagang(mahasiswa.ID, db) {
		return fmt.Errorf("kamu masih dalam proses seleksi/magang")
	}
	if mahasiswa.ID != idparam {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}

	namalengkap := r.FormValue("namalengkap")
	tanggallahir := r.FormValue("tanggallahir")
	jeniskelamin := r.FormValue("jeniskelamin")
	nim := r.FormValue("nim")
	perguruantinggi := r.FormValue("perguruantinggi")
	prodi := r.FormValue("prodi")

	if namalengkap == "" || tanggallahir == "" || jeniskelamin == "" || nim == "" || perguruantinggi == "" || prodi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	// Read the content of the file into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error 2: %s", err)
	}

	access_token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if access_token == "" {
		return fmt.Errorf("error access token: %s", err)
	}

	// Initialize GitHub client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access_token},
	)
	tc := oauth2.NewClient(r.Context(), ts)
	client := github.NewClient(tc)

	// Create a new repository file
	repoOwner := "Fatwaff"
	repoName := "bk-image"
	_, _, err = client.Repositories.CreateFile(r.Context(), repoOwner, repoName, "path/to/"+handler.Filename, &github.RepositoryContentFileOptions{
		Message:   github.String("Add new file"),
		Content:   fileContent,
		Committer: &github.CommitAuthor{Name: github.String("Fatwaff"), Email: github.String("fax.mp4@gmail.com")},
	})
	if err != nil {
		return fmt.Errorf("error 3: %s", err)
	}

	// imageUrl := "https://raw.githubusercontent.com/" + githubUser + "/" + repoName + "/master/" + imageFileName
	mhs := bson.M{
		"namalengkap":     namalengkap,
		"tanggallahir":    tanggallahir,
		"jeniskelamin":    jeniskelamin,
		"nim":             nim,
		"perguruantinggi": perguruantinggi,
		"prodi":           prodi,
		"seleksikampus":   0,
		"imagename":       "imageUrl",
		"akun": intermoni.User{
			ID: mahasiswa.Akun.ID,
		},
	}
	err = intermoni.UpdateOneDoc(idparam, db, "mahasiswa", mhs)
	if err != nil {
		return err
	}
	return nil
}

// by admin
func SeleksiMahasiswaByAdmin(idparam primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mahasiswa) error {
	mahasiswa, err := intermoni.GetMahasiswaFromID(idparam, db)
	if err != nil {
		return err
	}
	if CheckMahasiswa_MahasiswaMagang(mahasiswa.ID, db) {
		return fmt.Errorf("mahasiswa masih dalam proses seleksi")
	}
	if insertedDoc.SeleksiKampus != 1 && insertedDoc.SeleksiKampus != 2 {
		return fmt.Errorf("kesalahan server")
	}
	mahasiswa.SeleksiKampus = insertedDoc.SeleksiKampus
	err = intermoni.UpdateOneDoc(idparam, db, "mahasiswa", mahasiswa)
	if err != nil {
		return err
	}
	return nil
}

func CheckMahasiswa_MahasiswaMagang(idmahasiswa primitive.ObjectID, db *mongo.Database) bool {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mahasiswa._id": idmahasiswa,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	if count > 0 {
		jumlah_gagal := JumlahStatusGagalMahasiswaMagang(idmahasiswa, db)
		return jumlah_gagal != count
	}
	return false
}

func JumlahStatusGagalMahasiswaMagang(idmahasiswa primitive.ObjectID, db *mongo.Database) int64 {
	collection := db.Collection("mahasiswa_magang")
	filter := bson.M{
		"mahasiswa._id": idmahasiswa,
		"status": 2,
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0
	}
	return count
}