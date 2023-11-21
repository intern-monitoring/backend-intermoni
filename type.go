package intermoni

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Email				string				`bson:"email,omitempty" json:"email,omitempty"`
	Password			string				`bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword		string				`bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	Salt				string				`bson:"salt,omitempty" json:"salt,omitempty"`
	Role				string				`bson:"role,omitempty" json:"role,omitempty"`
}

type Password struct {
	Password			string				`bson:"password,omitempty" json:"password,omitempty"`
	Newpassword			string				`bson:"newpass,omitempty" json:"newpass,omitempty"`
	Confirmpassword		string				`bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
}

type Mahasiswa struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap			string				`bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	TanggalLahir		string				`bson:"tanggallahir,omitempty" json:"tanggallahir,omitempty"`
	JenisKelamin		string				`bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NIM					string				`bson:"nim,omitempty" json:"nim,omitempty"`
	PerguruanTinggi		string				`bson:"perguruantinggi,omitempty" json:"perguruantinggi,omitempty"`
	Prodi				string				`bson:"prodi,omitempty" json:"prodi,omitempty"`
	SeleksiKampus		int					`bson:"seleksikampus,omitempty" json:"seleksikampus,omitempty"`
	Akun				User				`bson:"akun,omitempty" json:"akun,omitempty"`
}

type Mitra struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	NamaNarahubung		string				`bson:"namanarahubung,omitempty" json:"namanarahubung,omitempty"`
	NoHpNarahubung		string				`bson:"nohpnarahubung,omitempty" json:"nohpnarahubung,omitempty"`
	Nama				string				`bson:"nama,omitempty" json:"nama,omitempty"`
	Kategori			string				`bson:"kategori,omitempty" json:"kategori,omitempty"`
	SektorIndustri		string				`bson:"sektorindustri,omitempty" json:"sektorindustri,omitempty"`
	Tentang				string				`bson:"tentang,omitempty" json:"tentang,omitempty"`
	Alamat				string				`bson:"alamat,omitempty" json:"alamat,omitempty"`
	Website				string				`bson:"website,omitempty" json:"website,omitempty"`
	MoU					int					`bson:"mou,omitempty" json:"mou,omitempty"`
	Akun				User				`bson:"akun,omitempty" json:"akun,omitempty"`
}

type Magang struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Posisi				string				`bson:"posisi,omitempty" json:"posisi,omitempty"`
	Mitra				Mitra				`bson:"mitra,omitempty" json:"mitra,omitempty"`
	Lokasi				string				`bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	CreatedAt			primitive.DateTime	`bson:"createdat,omitempty" json:"createdat,omitempty"`
	DeskripsiMagang		string				`bson:"deskripsimagang,omitempty" json:"deskripsimagang,omitempty"`
	InfoTambahanMagang	string				`bson:"infotambahanmagang,omitempty" json:"infotambahanmagang,omitempty"`
	Expired				string				`bson:"expired,omitempty" json:"expired,omitempty"`
}

type MahasiswaMagang struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Mahasiswa			Mahasiswa			`bson:"mahasiswa,omitempty" json:"mahasiswa,omitempty"`
	Magang				Magang				`bson:"magang,omitempty" json:"magang,omitempty"`
	Pembimbing			Pembimbing			`bson:"pembimbing,omitempty" json:"pembimbing,omitempty"`
	Mentor				Mentor				`bson:"mentor,omitempty" json:"mentor,omitempty"`
	SeleksiBerkas		int					`bson:"seleksiberkas,omitempty" json:"seleksiberkas,omitempty"`
	SeleksiWewancara	int					`bson:"seleksiwewancara,omitempty" json:"seleksiwewancara,omitempty"`
	Status				int					`bson:"status,omitempty" json:"status,omitempty"`
}

type Pembimbing struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap			string				`bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	NIK					string				`bson:"nik,omitempty" json:"nik,omitempty"`
	Prodi				string				`bson:"prodi,omitempty" json:"prodi,omitempty"`
	Akun				User				`bson:"akun,omitempty" json:"akun,omitempty"`
}

type Mentor struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap			string				`bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	NIK					string				`bson:"nik,omitempty" json:"nik,omitempty"`
	Mitra				Mitra				`bson:"mitra,omitempty" json:"mitra,omitempty"`
	Akun				User				`bson:"akun,omitempty" json:"akun,omitempty"`
}

type Report struct {
	ID					primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	MahasiswaMagang		MahasiswaMagang		`bson:"mahasiswamagang,omitempty" json:"mahasiswamagang,omitempty"`
	Judul				string				`bson:"judul,omitempty" json:"judul,omitempty"`
	Isi					string				`bson:"isi,omitempty" json:"isi,omitempty"`
	Penerima			User				`bson:"penerima,omitempty" json:"penerima,omitempty"`
	CreatedAt			primitive.DateTime	`bson:"createdat,omitempty" json:"createdat,omitempty"`
}

type Credential struct {
	Status				bool				`json:"status" bson:"status"`
	Token				string				`json:"token,omitempty" bson:"token,omitempty"`
	Message				string				`json:"message,omitempty" bson:"message,omitempty"`
	Role				string				`json:"role,omitempty" bson:"role,omitempty"`
}

type Response struct {
	Status				bool				`json:"status" bson:"status"`
	Message				string				`json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id					primitive.ObjectID	`json:"id"`
	Role				string				`json:"role"`
	Exp					time.Time			`json:"exp"`
	Iat					time.Time			`json:"iat"`
	Nbf					time.Time			`json:"nbf"`
}