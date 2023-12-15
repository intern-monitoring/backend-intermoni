package mahasiswa

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	githubUser     = "Fatwaff"
	repoName       = "bk-image"
	accessToken    = "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"
	uploadsDirPath = "uploads"
)

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

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	// Read the image data
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return fmt.Errorf("error 2: %s", err)
	}

	// Determine the file extension based on the actual image type
	ext := filepath.Ext(fileHeader.Filename)

	// Create the "uploads" directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return fmt.Errorf("error 3: %s", err)
	}

	// Create a new file in the server's storage with a dynamic file extension
	imageFileName := fmt.Sprintf("uploaded_image%s", ext)
	dst, err := os.Create(filepath.Join(uploadsDirPath, imageFileName))
	if err != nil {
		return fmt.Errorf("error 4: %s", err)
	}
	defer dst.Close()

	// Write the base64-encoded image data to the destination file
	_, err = dst.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error 5: %s", err)
	}

	// Upload the image to GitHub using GitHub API
	uploadURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubUser, repoName, imageFileName)
	encodedData := base64.StdEncoding.EncodeToString(buf.Bytes())
	reqBody := fmt.Sprintf(`{
		"message": "Upload image",
		"content": "%s"
	}`, encodedData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return fmt.Errorf("error 6: %s", err)
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error 7: %s", err)
	}
	defer resp.Body.Close()

	imageUrl := "https://raw.githubusercontent.com/" + githubUser + "/" + repoName + "/master/" + imageFileName
	mhs := bson.M{
		"namalengkap":     namalengkap,
		"tanggallahir":    tanggallahir,
		"jeniskelamin":    jeniskelamin,
		"nim":             nim,
		"perguruantinggi": perguruantinggi,
		"prodi":           prodi,
		"seleksikampus":   0,
		"imagename":       imageUrl,
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