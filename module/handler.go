package module

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/intern-monitoring/backend-intermoni/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// signup
func GCFHandlerSignUpMahasiswa(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var datamahasiswa model.Mahasiswa
	err := json.NewDecoder(r.Body).Decode(&datamahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpMahasiswa(conn, datamahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + datamahasiswa.NamaLengkap
	return GCFReturnStruct(Response)
}

func GCFHandlerSignUpMitra(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	var datamitra model.Mitra
	err := json.NewDecoder(r.Body).Decode(&datamitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUpMitra(conn, datamitra)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + datamitra.Nama
	return GCFReturnStruct(Response)
}

// login
func GCFHandlerLogin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, err := LogIn(conn, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := Encode(user.ID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email
		Response.Token = tokenstring
		Response.Role = user.Role
	}
	return GCFReturnStruct(Response)
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
	var Response model.Response
	Response.Status = false
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateEmailUser(payload.Id, conn, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Email"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdatePasswordUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	var Response model.Response
	Response.Status = false
	var datapassword model.Password
	err := json.NewDecoder(r.Body).Decode(&datapassword)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePasswordUser(payload.Id, conn, datapassword)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Password"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	var datauser model.User
	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateUser(payload.Id, conn, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update User"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllUserByAdmin(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "admin" {
		Response.Message = "Kamuh bukan admin"
		return GCFReturnStruct(Response)
	}
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetUser(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "admin" {
		return GCFHandlerGetUserFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllUserByAdmin(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetUserFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if data.Role == "mahasiswa" {
		datamahasiswa, err := GetMahasiswaFromAkun(data.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		datamahasiswa.Akun = data
		return GCFReturnStruct(datamahasiswa) 
	}
	if data.Role == "mitra" {
		datamitra, err := GetMitraFromAkun(data.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		datamitra.Akun = data
		return GCFReturnStruct(datamitra) 
	}
	Response.Message = "Tidak ada data"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetUserFromID(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	data, err := GetUserFromID(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

// mahasiswa
func GCFHandlerUpdateMahasiswa(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mahasiswa" {
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
	var datamahasiswa model.Mahasiswa
	err = json.NewDecoder(r.Body).Decode(&datamahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMahasiswa(idparam, payload.Id, conn, datamahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Mahasiswa"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMahasiswa(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
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
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mahasiswa" {
		Response.Message = "Maneh bukan mahasiswa"
		return GCFReturnStruct(Response)
	}
	data, err := GetMahasiswaFromAkun(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

// mitra
func GCFHandlerUpdateMitra(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
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
	var datamitra model.Mitra
	err = json.NewDecoder(r.Body).Decode(&datamitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMitra(idparam, payload.Id, conn, datamitra)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Mitra"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMitra(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
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
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
		Response.Message = "Maneh bukan mitra"
		return GCFReturnStruct(Response)
	}
	data, err := GetMitraFromAkun(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

// magang
func GCFHandlerInsertMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	var datamagang model.Magang
	err = json.NewDecoder(r.Body).Decode(&datamagang)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = InsertMagang(payload.Id, conn, datamagang)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Insert Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
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
	var datamagang model.Magang
	err = json.NewDecoder(r.Body).Decode(&datamagang)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdateMagang(idparam, payload.Id, conn, datamagang)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Update Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
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
	err = DeleteMagang(idparam, payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Delete Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllMagang(MONGOCONNSTRINGENV, dbname string) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	data, err := GetAllMagang(conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMagangFromID(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetAllMagang(MONGOCONNSTRINGENV, dbname)
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetMagangFromID(objID, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMagangFromMitra(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
		Response.Message = "Kamuh bukan Mitra"
		return GCFReturnStruct(Response)
	}
	data, err := GetMagangFromMitra(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

func GCFHandlerGetMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mitra" {
		return GCFHandlerGetMagangFromID(MONGOCONNSTRINGENV, dbname, r)
	}
	id := GetID(r)
	if id == "" {
		return GCFHandlerGetMagangFromMitra(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname, r)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}
	data, err := GetMagangFromIDByMitra(idparam, payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	return GCFReturnStruct(data)
}

// mahasiswa magang
func GCFHandlerInsertMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	mahasiswa, err := GetMahasiswaFromAkun(payload.Id, conn)
	if err != nil {
		Response.Message = err.Error()
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
	err = InsertMahasiswaMagang(idmagang, mahasiswa.ID, conn)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Apply Magang"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)	
	var Response model.Response
	Response.Status = false
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return GCFReturnStruct(Response)
	}
	if payload.Role != "admin" && payload.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return GCFReturnStruct(Response)
	}
	if payload.Role == "admin" {
		data, err := GetMahasiswaMagangByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(data)
	}
	if payload.Role == "mitra" {
		data, err := GetMahasiswaMagangByMitra(payload.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}
		return GCFReturnStruct(data)
	}
	Response.Message = "Tidak ada data"
	return GCFReturnStruct(Response)
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// get id
func GetID(r *http.Request) string {
    return r.URL.Query().Get("id")
}