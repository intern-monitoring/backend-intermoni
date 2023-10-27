package Fatwa_1214038

import (
	"fmt"
	"testing"

	"github.com/intern-monitoring/backend-intermoni/model"
	"github.com/intern-monitoring/backend-intermoni/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "db_intermoni")

func TestGetUserFromEmail(t *testing.T) {
	email := "admin@gmail.com"
	hasil, err := module.GetUserFromEmail(email, db)
	if err != nil {
		t.Errorf("Error TestGetUserFromEmail: %v", err)
	} else {
		fmt.Println(hasil)
	}
}

func TestInsertOneMagang(t *testing.T) {
	var doc model.Magang
//    doc.Logo = "https://fatwaff.github.io/bk-image/user/ford.jpg"
   doc.Posisi = "Network Engineer"
   doc.Perusahaan = "Ford Company Etc"
   doc.Lokasi = "Bandung"
   doc.CreatedAt = "07-08-2023"
   doc.DeskripsiMagang = "<div><ul><li>Mengurus administrasi bagian marketing</li><li>Membuat Sales Order,membuat Penawaran Harga</li><li>Menerima Purchase Order (PO) Customer</li><li>Membina hubungan baik antara Perusahaan dan Customer</li><li>Bisa bekerja secara akurat dan memperhatikan detail sehingga bisa memproses pesanan dengan cepat dan efisien</li><li>Jujur, pekerja keras,ulet,tekun,bertanggung jawab,punya komitmen yang tinggi, percaya diri, memiliki kemampuan komunikasi yang baik</li></ul></div>"
   doc.InfoTambahanMagang = "<div><ul><li>Pengalaman 3 tahun kerja</li><li>Pegawai tetap</li></ul></div>"
   doc.TentangPerusahaan = "<div>Ford Motor Company (commonly known as Ford) is an American multinational automobile manufacturer headquartered in Dearborn, Michigan, United States. It was founded by Henry Ford and incorporated on June 16, 1903. The company sells automobiles and commercial vehicles under the Ford brand, and luxury cars under its Lincoln brand.</div>"
//    doc.InfoTambahanPerusahaan = "<div><ul><li>1000-2000 Pekerja</li><li>Industri Manufaktur/Produksi</li></ul></div>"
   if  doc.Posisi == "" || doc.Perusahaan == "" || doc.Lokasi == "" || doc.CreatedAt == "" || doc.DeskripsiMagang == "" || doc.InfoTambahanMagang == "" || doc.TentangPerusahaan == "" {
	   t.Errorf("mohon untuk melengkapi data")
   } else {
	   insertedID, err := module.InsertOneDoc(db, "lowongan", doc)
	   if err != nil {
		   t.Errorf("Error inserting document: %v", err)
		   fmt.Println("Data tidak berhasil disimpan")
	   } else {
	   fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
	   }
   }
}

type Userr struct {
	ID           	primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email  			string             `bson:"email,omitempty" json:"email,omitempty"`
	Role     		string			   `bson:"role,omitempty" json:"role,omitempty"`
}

func TestGetAllDoc(t *testing.T) {
	hasil := module.GetAllDocs(db, "user", []Userr{})
	fmt.Println(hasil)
}

// func TestUpdateOneDoc(t *testing.T) {
//  	var docs model.User
// 	id := "649063d3ad72e074286c61e8"
// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	docs.FirstName = "Aufah"
// 	docs.LastName = "Auliana"
// 	docs.Email = "aufa@gmail.com"
// 	docs.Password = "123456"
// 	if docs.FirstName == "" || docs.LastName == "" || docs.Email == "" || docs.Password == "" {
// 		t.Errorf("mohon untuk melengkapi data")
// 	} else {
// 		err := module.UpdateOneDoc(db, "user", objectId, docs)
// 		if err != nil {
// 			t.Errorf("Error inserting document: %v", err)
// 			fmt.Println("Data tidak berhasil diupdate")
// 		} else {
// 			fmt.Println("Data berhasil diupdate")
// 		}
// 	}
// }

// func TestGetLowonganFromID(t *testing.T){
// 	id := "64d0b1104255ba95ba588512"
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil{
// 		t.Fatalf("error converting id to objectID: %v", err)
// 	}
// 	doc, err := module.GetLowonganFromID(objectId)
// 	if err != nil {
// 		t.Fatalf("error calling GetDocFromID: %v", err)
// 	}
// 	fmt.Println(doc)
// }

func TestSignUpMahasiswa(t *testing.T) {
	var doc model.Mahasiswa
	doc.NamaLengkap = "Erdito Nausha Adam"
	doc.TanggalLahir = "20/05/2001"
	doc.JenisKelamin = "Laki-laki"
	doc.NIM = "1214031"
	doc.PerguruanTinggi = "Universitas Logistik dan Bisnis Internasional"
	doc.Prodi = "D4 Teknik Informatika"
	doc.Akun.Email = "erdito2@gmail.com"
	doc.Akun.Password = "fghjkliow"
	doc.Akun.Confirmpassword = "fghjkliow"
	err := module.SignUpMahasiswa(db, "mahasiswa", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaLengkap)
	}
}

func TestSignUpMitra(t *testing.T) {
	var doc model.Mitra
	doc.NamaNarahubung = "Erdito Nausha Adam"
	doc.NoHpNarahubung = "085728980009"
	doc.NamaResmi = "PT. Maju Mundur"
	doc.Kategori = "BUMN"
	doc.SektorIndustri = "Teknologi Informasi"
	doc.Alamat = "Jl. Sariasih 2"
	doc.Website = "www.majumundur.com"
	doc.Akun.Email = "erdito@gmail.com"
	doc.Akun.Password = "fghjkliow"
	doc.Akun.Confirmpassword = "fghjkliow"
	err := module.SignUpMitra(db, "mitra", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaResmi)
	}
}

// func TestSignUpIndustri(t *testing.T) {
// 	var doc model.User
// 	doc.FirstName = "Dimas"
// 	doc.LastName = "Ardianto"
// 	doc.Email = "dimas@gmail.com"
// 	doc.Password = "fghjkliow"
// 	doc.Confirmpassword = "fghjkliow"
// 	insertedID, err := module.SignUpIndustri(db, "user", doc)
// 	if err != nil {
// 		t.Errorf("Error inserting document: %v", err)
// 	} else {
// 	fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
// 	}
// }

func TestLogIn(t *testing.T) {
	var doc model.User
	doc.Email = "dimas@gmail.com"
	doc.Password = "fghjkliow"
	user, err := module.LogIn(db, doc)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome bang:", user)
	}
}

func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := module.GenerateKey()
	fmt.Println("ini private key :", privateKey)
	fmt.Println("ini public key :", publicKey)
	id := "64d0b1104255ba95ba588512"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	hasil, err := module.Encode(objectId, privateKey)
	fmt.Println("ini hasil :", hasil, err)
}

// func TestWatoken2(t *testing.T) {
// 	var user model.User
// 	privateKey, publicKey := module.GenerateKey()
// 	fmt.Println("privateKey : ", privateKey)
// 	fmt.Println("publicKey : ", publicKey)
// 	id := "649063d3ad72e074286c61e8"
// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	user.FirstName = "Fatwa"
// 	user.LastName = "Fatahillah"
// 	user.Email = "fax@gmail.com"
// 	user.Role = "pelamar"
// 	tokenstring, err := module.Encode(objectId, privateKey)
// 	if err != nil {
// 		t.Errorf("Error getting document: %v", err)
// 	} else {
// 		body, err := module.Decode(publicKey, tokenstring)
// 		fmt.Println("signed : ", tokenstring)
// 		fmt.Println("isi : ", body)
// 		if err != nil {
// 			t.Errorf("Error getting document: %v", err)
// 		} else {
// 			fmt.Println("Berhasil yey!")
// 		}
// 	}
// }

func TestWatoken(t *testing.T) {
	body, err := module.Decode("f3248b509d9731ebd4e0ccddadb5a08742e083db01678e8a1d734ce81298868f", "v4.public.eyJlbWFpbCI6ImZheEBnbWFpbC5jb20iLCJleHAiOiIyMDIzLTEwLTIyVDAwOjQxOjQ1KzA3OjAwIiwiZmlyc3RuYW1lIjoiRmF0d2EiLCJpYXQiOiIyMDIzLTEwLTIxVDIyOjQxOjQ1KzA3OjAwIiwiaWQiOiI2NDkwNjNkM2FkNzJlMDc0Mjg2YzYxZTgiLCJsYXN0bmFtZSI6IkZhdGFoaWxsYWgiLCJuYmYiOiIyMDIzLTEwLTIxVDIyOjQxOjQ1KzA3OjAwIiwicm9sZSI6InBlbGFtYXIifR_Q4b9X7WC7up7dUUxz_Yki39M-ReovTIoTFfdJmFYRF5Oer0zQZx_ZQamkOsogJ6RuGJhxT3OxrXFS5p6dMg0")
	fmt.Println("isi : ", body, err)
}

// func TestWatoken3(t *testing.T) {
// 	var datauser model.User
// 	privateKey, publicKey := module.GenerateKey()
// 	fmt.Println("privateKey : ", privateKey)
// 	fmt.Println("publicKey : ", publicKey)
// 	datauser.Email = "fatwaff@gmail.com"
// 	datauser.Password = "fghjklio"
// 	user, err := module.LogIn(db, "user", datauser)
// 	fmt.Println("id : ", user.ID)
// 	fmt.Println("firstname : ", user.FirstName)
// 	fmt.Println("lastname : ", user.LastName)
// 	fmt.Println("email : ", user.Email)
// 	fmt.Println("role : ", user.Role)
// 	if err != nil {
// 		t.Errorf("Error getting document: %v", err)
// 	} else {
// 		tokenstring, err := module.Encode(user.ID, privateKey)
// 		if err != nil {
// 			t.Errorf("Error getting document: %v", err)
// 		} else {
// 			body, err := module.Decode(publicKey, tokenstring)
// 			fmt.Println("signed : ", tokenstring)
// 			fmt.Println("isi : ", body)
// 			if err != nil {
// 				t.Errorf("Error getting document: %v", err)
// 			} else {
// 				fmt.Println("Berhasil yey!")
// 			}
// 		}
// 	}
// }