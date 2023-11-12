package Fatwa_1214038

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/intern-monitoring/backend-intermoni/model"
	"github.com/intern-monitoring/backend-intermoni/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
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
   doc.Lokasi = "Bandung"
   doc.DeskripsiMagang = "<div><ul><li>Mengurus administrasi bagian marketing</li><li>Membuat Sales Order,membuat Penawaran Harga</li><li>Menerima Purchase Order (PO) Customer</li><li>Membina hubungan baik antara Mitra dan Customer</li><li>Bisa bekerja secara akurat dan memperhatikan detail sehingga bisa memproses pesanan dengan cepat dan efisien</li><li>Jujur, pekerja keras,ulet,tekun,bertanggung jawab,punya komitmen yang tinggi, percaya diri, memiliki kemampuan komunikasi yang baik</li></ul></div>"
   doc.InfoTambahanMagang = "<div><ul><li>Pengalaman 3 tahun kerja</li><li>Pegawai tetap</li></ul></div>"
   if  doc.Posisi == "" || doc.Lokasi == "" || doc.DeskripsiMagang == "" || doc.InfoTambahanMagang == "" {
	   t.Errorf("mohon untuk melengkapi data")
   } else {
	   insertedID, err := module.InsertOneDoc(db, "magang", doc)
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

func TestInsertUser(t *testing.T) {
	var doc model.User
	doc.Email = "admin@gmail.com"
	password := "admin123"
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		t.Errorf("kesalahan server : salt")
	} else {
		hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
		user := bson.M{
			"email": doc.Email,
			"password": hex.EncodeToString(hashedPassword),
			"salt": hex.EncodeToString(salt),
			"role": "admin",
		}
		_, err = module.InsertOneDoc(db, "user", user)
		if err != nil {
			t.Errorf("gagal insert")
		} else {
			fmt.Println("berhasil insert")
		}
	}
}

func TestGetUserByAdmin(t *testing.T) {
	id := "65473763d04dda3a8502b58f"
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to objectID: %v", err)
	}
	data, err := module.GetUserFromID(idparam, db)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		if data.Role == "mahasiswa" {
			datamahasiswa, err := module.GetMahasiswaFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				datamahasiswa.Akun = data
				fmt.Println(datamahasiswa) 
			}
		}
		if data.Role == "mitra" {
			datamitra, err := module.GetMitraFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				datamitra.Akun = data
				fmt.Println(datamitra)
			}
		}
	}
}

func TestSignUpMahasiswa(t *testing.T) {
	var doc model.Mahasiswa
	doc.NamaLengkap = "Adit Nausha Adam"
	doc.TanggalLahir = "20/05/2001"
	doc.JenisKelamin = "Laki-laki"
	doc.NIM = "1214031"
	doc.PerguruanTinggi = "Universitas Logistik dan Bisnis Internasional"
	doc.Prodi = "D4 Teknik Informatika"
	doc.Akun.Email = "adit@gmail.com"
	doc.Akun.Password = "fghjkliow"
	doc.Akun.Confirmpassword = "fghjkliow"
	err := module.SignUpMahasiswa(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaLengkap)
	}
}

func TestSignUpMitra(t *testing.T) {
	var doc model.Mitra
	doc.NamaNarahubung = "Dimas Ardianto"
	doc.NoHpNarahubung = "085728980009"
	doc.Nama = "PT. Mundur Maju"
	doc.Kategori = "BUMN"
	doc.SektorIndustri = "Teknologi Informasi"
	doc.Alamat = "Jl. Sariasih 3"
	doc.Website = "www.mundurmaju.com"
	doc.Akun.Email = "dimas@gmail.com"
	doc.Akun.Password = "fghjkliow"
	doc.Akun.Confirmpassword = "fghjkliow"
	err := module.SignUpMitra(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
	fmt.Println("Data berhasil disimpan dengan nama :", doc.Nama)
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
	role := "admin"
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	hasil, err := module.Encode(objectId, role, privateKey)
	fmt.Println("ini hasil :", hasil, err)
}

func TestUpdateMahasiswa(t *testing.T) {
	var doc model.Mahasiswa
	id := "654a01bde89e6f232a62e41d"
	objectId, _ := primitive.ObjectIDFromHex(id)
	id2 := "654a01bce89e6f232a62e41b"
	userid, _ := primitive.ObjectIDFromHex(id2)
	doc.NamaLengkap = "Adito Nausha Adam"
	doc.TanggalLahir = "20/05/2001"
	doc.JenisKelamin = "Laki-laki"
	doc.NIM = "1214031"
	doc.PerguruanTinggi = "Universitas Logistik dan Bisnis Internasional"
	if doc.NamaLengkap == "" || doc.TanggalLahir == "" || doc.JenisKelamin == "" || doc.NIM == "" || doc.PerguruanTinggi == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		err := module.UpdateMahasiswa(objectId, userid, db, doc)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil diupdate")
		} else {
			fmt.Println("Data berhasil diupdate")
		}
	}
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


// test magang
func TestInsertMagang(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTEwLTMxVDA4OjQ4OjIyWiIsImlhdCI6IjIwMjMtMTAtMzFUMDY6NDg6MjJaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTEwLTMxVDA2OjQ4OjIyWiJ9lXy1b5nOEYuCn7_o-TcFuR-3OOm__T7SHlAdx3PQl4Du9EAr8pu85lvU6SddRar7YB3DEbf-zwfY_zytj7jrAQ")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	var datamagang model.Magang
	datamagang.Posisi = "Data Science"
	datamagang.Lokasi = "Bandung"
	datamagang.DeskripsiMagang = "<div><ul><li>Mengurus administrasi bagian marketing</li><li>Membuat Sales Order,membuat Penawaran Harga</li><li>Menerima Purchase Order (PO) Customer</li><li>Membina hubungan baik antara Mitra dan Customer</li><li>Bisa bekerja secara akurat dan memperhatikan detail sehingga bisa memproses pesanan dengan cepat dan efisien</li><li>Jujur, pekerja keras,ulet,tekun,bertanggung jawab,punya komitmen yang tinggi, percaya diri, memiliki kemampuan komunikasi yang baik</li></ul></div>"
	datamagang.InfoTambahanMagang = "<div><ul><li>Pengalaman 3 tahun kerja</li><li>Pegawai tetap</li></ul></div>"
	datamagang.Expired = "01-11-2023"
	err = module.InsertMagang(payload.Id, conn, datamagang)
	if err != nil {
		t.Errorf("Error insert : %v", err)
	} else {
		fmt.Println("Berhasil yey!")
	}
}

func TestUpdateMagang(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTExLTAxVDA2OjQ5OjQ0WiIsImlhdCI6IjIwMjMtMTEtMDFUMDQ6NDk6NDRaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTExLTAxVDA0OjQ5OjQ0WiJ92RxBGslXaHBoLQhvMJLQO7uEBG5c5FmkpZgakPjmk1aUFDdRkw3m3r-7BpkhDmCtByoARDr36X3DhjCL8HT8AQ")
	// payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTExLTAxVDA2OjQ3OjMxWiIsImlhdCI6IjIwMjMtMTEtMDFUMDQ6NDc6MzFaIiwiaWQiOiI2NTNkZTllYjg5MzlmYjNjZjI3ZjZkMzciLCJuYmYiOiIyMDIzLTExLTAxVDA0OjQ3OjMxWiJ92YbTLQWznLupbH0Syb6GPKkj4ZW_JJLveVcFTfZElv8_jybZCMBnw8y-7SLZVMpRTq56PaArdEBwlvvSPQjtCg")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	var datamagang model.Magang
	datamagang.Posisi = "Data Sciences"
	datamagang.Lokasi = "Bandung"
	datamagang.DeskripsiMagang = "<div><ul><li>Mengurus administrasi bagian marketing</li><li>Membuat Sales Order,membuat Penawaran Harga</li><li>Menerima Purchase Order (PO) Customer</li><li>Membina hubungan baik antara Mitra dan Customer</li><li>Bisa bekerja secara akurat dan memperhatikan detail sehingga bisa memproses pesanan dengan cepat dan efisien</li><li>Jujur, pekerja keras,ulet,tekun,bertanggung jawab,punya komitmen yang tinggi, percaya diri, memiliki kemampuan komunikasi yang baik</li></ul></div>"
	datamagang.InfoTambahanMagang = "<div><ul><li>Pengalaman 3 tahun kerja</li><li>Pegawai tetap</li></ul></div>"
	datamagang.Expired = "01-11-2023"
	id := "65406377996edfaee3ed9a19"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	err = module.UpdateMagang(objectId, payload.Id, conn, datamagang)
	if err != nil {
		t.Errorf("Error update : %v", err)
	} else {
		fmt.Println("Berhasil yey!")
	}
}

func TestDeleteMagang(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTExLTAxVDA2OjQ5OjQ0WiIsImlhdCI6IjIwMjMtMTEtMDFUMDQ6NDk6NDRaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTExLTAxVDA0OjQ5OjQ0WiJ92RxBGslXaHBoLQhvMJLQO7uEBG5c5FmkpZgakPjmk1aUFDdRkw3m3r-7BpkhDmCtByoARDr36X3DhjCL8HT8AQ")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	id := "65406377996edfaee3ed9a19"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	err = module.DeleteMagang(objectId, payload.Id, conn)
	if err != nil {
		t.Errorf("Error delete : %v", err)
	} else {
		fmt.Println("Berhasil yey!")
	}
}

func TestGetMagangByMitra(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTEwLTMxVDEwOjIxOjI4WiIsImlhdCI6IjIwMjMtMTAtMzFUMDg6MjE6MjhaIiwiaWQiOiI2NTQwNjMyODI4NzY0ZDk2YzY0OWYyOWQiLCJuYmYiOiIyMDIzLTEwLTMxVDA4OjIxOjI4WiJ9CoWv7X_t-idleCPyTX3jvwbcSR038WX6av6gmh8hpAV5_B5Moe11GK-hpz-osTdzpAuTUw0Qsueic9ny0qg1DQ")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	magang, err := module.GetMagangFromMitra(payload.Id, conn)
	if err != nil {
		t.Errorf("Error get magang : %v", err)
	} else {
		fmt.Println(magang)
	}
}

func TestGetAllMagang(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	data, err := module.GetAllMagang(conn)
	if err != nil {
		t.Errorf("Error get all : %v", err)
	} else {
		fmt.Println(data)
	}
}

func TestGetMagangFromID(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	id := "65406377996edfaee3ed9a19"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	magang, err := module.GetMagangFromID(objectId, conn)
	if err != nil {
		t.Errorf("Error get magang : %v", err)
	} else {
		fmt.Println(magang)
	}
}

func TestGetMagangFromIDByMitra(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	payload, err := module.Decode("b95509d9634ed137b5ccdd07a7534ab0dcede0f310c09634afbf0262c7a4ce1c", "v4.public.eyJleHAiOiIyMDIzLTEwLTMxVDE0OjI1OjQxWiIsImlhdCI6IjIwMjMtMTAtMzFUMTI6MjU6NDFaIiwiaWQiOiI2NTNkZTllYjg5MzlmYjNjZjI3ZjZkMzciLCJuYmYiOiIyMDIzLTEwLTMxVDEyOjI1OjQxWiJ9RUHYj4xe2MmJcABLLiQ_ftjiNiM2CW6ABhSY0ovQ9SL0uQ9AK2m2v7svW2LU5u8XWB4tQUjQchelIId8KzITDQ")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	// if payload.Role != "mitra" {
	// 	t.Errorf("Error role: %v", err)
	// }
	id := "654060d83f526c35452232cf"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		t.Fatalf("error converting id to objectID: %v", err)
	}
	magang, err := module.GetMagangFromIDByMitra(objectId, payload.Id, conn)
	if err != nil {
		t.Errorf("Error get magang : %v", err)
	} else {
		fmt.Println(magang)
	}
}

func TestInsertMahasiswaMagang(t *testing.T) {
	idmg := "654a64e04c78636f1e657479"
	idmhs := "65473763d04dda3a8502b58f"
	idmgid, _ := primitive.ObjectIDFromHex(idmg)
	idmhsid, _ := primitive.ObjectIDFromHex(idmhs)
	err := module.InsertMahasiswaMagang(idmgid, idmhsid, db)
	if err != nil {
		t.Errorf("Error insert : %v", err)
	} else {
		fmt.Println("Berhasil yey!")
	}
}

func TestSeleksiMahasiswaMagangByAdmin(t *testing.T) {
	idmg := "6550cd47b25dc17957071a4a"
	idmgid, _ := primitive.ObjectIDFromHex(idmg)
	err := module.SeleksiMahasiswaMagangByAdmin(idmgid, db)
	if err != nil {
		t.Errorf("Error get mahasiswa : %v", err)
	} else {
		fmt.Println("Berhasil yey!")
	}
}

func TestGetMahasiswaMagangByAdmin(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "db_intermoni")
	mahasiswa, err := module.GetMahasiswaMagangByAdmin(conn)
	if err != nil {
		t.Errorf("Error get mahasiswa : %v", err)
	} else {
		fmt.Println(mahasiswa)
	}
}

func TestReturnStruct(t *testing.T){
	// var user model.User
	// user.Email = "fatwa"
	id := "65473763d04dda3a8502b58f"
	objectId, _ := primitive.ObjectIDFromHex(id)
	user, _ := module.GetUserFromID(objectId, db)
	data := model.User{ 
		ID : user.ID,
		Email: user.Email,
		Role : user.Role,
	}
	hasil := module.GCFReturnStruct(data)
	fmt.Println(hasil)
}