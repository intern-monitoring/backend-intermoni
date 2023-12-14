package mahasiswa

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
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
func UpdateMahasiswa(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.Mahasiswa) error {
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
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NIM == "" || insertedDoc.PerguruanTinggi == "" || insertedDoc.Prodi == "" || insertedDoc.Image == nil  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	// Access the FormData from the JSON
	imageData := insertedDoc.Image

	// Detect the content type of the image
	contentType := http.DetectContentType(imageData)

	// Determine the file extension based on content type
	fileExt, err := extFromContentType(contentType)
	if err != nil {
		return err
	}

	// Create a new file in the server's storage
	imageFileName := "uploaded_image" + fileExt
	dst, err := os.Create(filepath.Join(uploadsDirPath, imageFileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Write the base64-encoded image data to the destination file
	_, err = dst.Write(imageData)
	if err != nil {
		return err
	}

	// Upload the image to GitHub using GitHub API
	uploadURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubUser, repoName, imageFileName)
	reqBody := fmt.Sprintf(`{
		"message": "Upload image",
		"content": "%s"
	}`, base64.StdEncoding.EncodeToString(imageData))

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	imageUrl := "https://raw.githubusercontent.com/" + githubUser + "/" + repoName + "/master/" + imageFileName
	mhs := bson.M{
		"namalengkap":     insertedDoc.NamaLengkap,
		"tanggallahir":    insertedDoc.TanggalLahir,
		"jeniskelamin":    insertedDoc.JenisKelamin,
		"nim":             insertedDoc.NIM,
		"perguruantinggi": insertedDoc.PerguruanTinggi,
		"prodi":           insertedDoc.Prodi,
		"seleksikampus":   0,
		"imageurl":        imageUrl,
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

// extFromContentType returns the file extension based on the given content type.
func extFromContentType(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	// Add more cases for other image types as needed
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
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