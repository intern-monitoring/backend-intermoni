package module

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/intern-monitoring/backend-intermoni/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response model.Response
	user model.User
	mahasiswa model.Mahasiswa
	mitra model.Mitra
	magang model.Magang
)

// signup
func GCFHandlerSignUpMahasiswa(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpMahasiswa(conn, mahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Halo " + mahasiswa.NamaLengkap
	return GCFReturnStruct(Response)
}

func GCFHandlerSignUpMitra(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&mitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpMitra(conn, mitra)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Halo " + mitra.Nama
	return GCFReturnStruct(Response)
}

// login
func GCFHandlerLogin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Credential model.Credential
	Credential.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Credential)
	}
	user, err := LogIn(conn, user)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Credential)
	}
	tokenstring, err := Encode(user.ID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Credential.Message = "Gagal Encode Token : " + err.Error()
		return GCFReturnStruct(Credential)
	}
	//
	Credential.Message = "Selamat Datang " + user.Email
	Credential.Token = tokenstring
	Credential.Role = user.Role
	Credential.Status = true
	return GCFReturnStruct(Credential)
}

// get all
func GCFHandlerGetAll(MONGOCONNSTRINGENV, dbname, col string, docs interface{}) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	data := GetAllDocs(conn, col, docs)
	return GCFReturnStruct(data)
}

// user
func GCFHandlerUpdateEmailUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateEmailUser(user_login.Id, conn, user)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Email"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdatePasswordUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var password model.Password
	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePasswordUser(user_login.Id, conn, password)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Password"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateUser(user_login.Id, conn, user)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update User"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "admin" {
		return GCFHandlerGetUserFromID(user_login.Id, conn, r)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllUserByAdmin(conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	user, err := GetUserFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user.Role == "mahasiswa" {
		mahasiswa, err := GetMahasiswaFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		mahasiswa.Akun = user
		return GCFReturnStruct(mahasiswa) 
	}
	if user.Role == "mitra" {
		mitra, err := GetMitraFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		mitra.Akun = user
		return GCFReturnStruct(mitra) 
	}
	//
	Response.Message = "Tidak ada data"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetUserFromID(iduser primitive.ObjectID, conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	user, err := GetUserFromID(iduser, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	result := bson.M{
		"_id" : user.ID,
		"email": user.Email,
		"role" : user.Role,
	}
	return GCFReturnStruct(result)
}

func GCFHandlerGetAllUserByAdmin(conn *mongo.Database) string {
	Response.Status = false
	//
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(data)
}

// mahasiswa
func GCFHandlerUpdateMahasiswa(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMahasiswa(idparam, user_login.Id, conn, mahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mahasiswa"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMahasiswa(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	
	data, err := GetAllMahasiswa(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMahasiswaFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh bukan mahasiswa"
		return GCFReturnStruct(Response)
	}
	mahasiswa, err := GetMahasiswaFromAkun(user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(mahasiswa)
}

// mitra
func GCFHandlerUpdateMitra(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMitra(idparam, user_login.Id, conn, mitra)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mitra"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMitra(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	
	data, err := GetAllMitra(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMitraFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh bukan mitra"
		return GCFReturnStruct(Response)
	}
	mitra, err := GetMitraFromAkun(user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(mitra)
}

// magang
func GCFHandlerInsertMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&magang)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = InsertMagang(user_login.Id, conn, magang)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Insert Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&magang)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMagang(idparam, user_login.Id, conn, magang)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = DeleteMagang(idparam, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Delete Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMagang(conn *mongo.Database) string {
	Response.Status = false
	
	data, err := GetAllMagang(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMagangFromID(conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllMagang(conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	magang, err := GetMagangFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(magang)
}

func GCFHandlerGetMagangFromMitra(idmitra primitive.ObjectID, conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	magang, err := GetMagangFromMitra(idmitra, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(magang)
}

func GCFHandlerGetMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		return GCFHandlerGetMagangFromID(conn, r)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetMagangFromMitra(user_login.Id, conn, r)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	magang, err := GetMagangFromIDByMitra(idparam, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(magang)
}

// mahasiswa magang
func GCFHandlerInsertMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idmagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = InsertMahasiswaMagang(idmagang, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Apply Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idmahasiswamagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	err = DeleteMahasiswaMagang(idmahasiswamagang, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Batal Apply"
	return GCFReturnStruct(Response)

}

func GCFHandlerSeleksiMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idmahasiswamagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	if user_login.Role == "admin" {
		err = SeleksiMahasiswaMagangByAdmin(idmahasiswamagang, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		Response.Status = true
		Response.Message = "Berhasil Seleksi"
	}
	if user_login.Role == "mitra" {
		err = SeleksiMahasiswaMagangByMitra(idmahasiswamagang, user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		Response.Status = true
		Response.Message = "Berhasil Seleksi"
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	Response.Status = false
	//
	user_login, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllMahasiswaMagang(user_login, conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	mahasiswa_magang, err := GetDetailMahasiswaMagangFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	//
	return GCFReturnStruct(mahasiswa_magang)
}

func GCFHandlerGetAllMahasiswaMagang(user_login model.Payload, conn *mongo.Database) string {
	Response.Status = false
	//
	if user_login.Role == "admin" {
		mahasiswa_magang, err := GetMahasiswaMagangByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(mahasiswa_magang)
	}
	if user_login.Role == "mitra" {
		mahasiswa_magang, err := GetMahasiswaMagangByMitra(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(mahasiswa_magang)
	}
	if user_login.Role == "mahasiswa" {
		mahasiswa_magang, err := GetMahasiswaMagangByMahasiswa(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(mahasiswa_magang)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return GCFReturnStruct(Response)
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (model.Payload, error) {
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