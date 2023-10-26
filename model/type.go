package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           	primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email  			string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string         	   `bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword string         	   `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	Salt 			string			   `bson:"salt,omitempty" json:"salt,omitempty"`
	Role     		string			   `bson:"role,omitempty" json:"role,omitempty"`
}

type Mahasiswa struct {
	ID           	primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  	string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	TanggalLahir	string             `bson:"tanggallahir,omitempty" json:"tanggallahir,omitempty"`
	JenisKelamin  	string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NIM  			string             `bson:"nim,omitempty" json:"nim,omitempty"`
	PerguruanTinggi string             `bson:"perguruantinggi,omitempty" json:"perguruantinggi,omitempty"`
	Prodi  			string             `bson:"prodi,omitempty" json:"prodi,omitempty"`
	Akun     		User			   `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Mitra struct {
	ID           	primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaNarahubung  string             `bson:"namanarahubung,omitempty" json:"namanarahubung,omitempty"`
	NoHpNarahubung  string             `bson:"nohpnarahubung,omitempty" json:"nohpnarahubung,omitempty"`
	NamaResmi  		string             `bson:"namaresmi,omitempty" json:"namaresmi,omitempty"`
	Kategori 		string             `bson:"kategori,omitempty" json:"kategori,omitempty"`
	SektorIndustri 	string             `bson:"sektorindustri,omitempty" json:"sektorindustri,omitempty"`
	Alamat 			string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Website 		string             `bson:"website,omitempty" json:"website,omitempty"`
	Akun     		User			   `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Magang struct {
	ID              		primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Posisi         			string			   `bson:"posisi,omitempty" json:"posisi,omitempty"`
	Perusahaan      		string			   `bson:"perusahaan,omitempty" json:"perusahaan,omitempty"`
	Lokasi          		string			   `bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	CreatedAt       		string			   `bson:"createdat,omitempty" json:"createdat,omitempty"`
	DeskripsiMagang  		string			   `bson:"deskripsimagang,omitempty" json:"deskripsimagang,omitempty"`
	InfoTambahanMagang   	string			   `bson:"infotambahanmagang,omitempty" json:"infotambahanmagang,omitempty"`
	TentangPerusahaan   	string			   `bson:"tentangperusahaan,omitempty" json:"tentangperusahaan,omitempty"`
	Expired   				string			   `bson:"expired,omitempty" json:"expired,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id           	primitive.ObjectID `json:"id"`
	Exp 			time.Time 	 	   `json:"exp"`
	Iat 			time.Time 		   `json:"iat"`
	Nbf 			time.Time 		   `json:"nbf"`
}