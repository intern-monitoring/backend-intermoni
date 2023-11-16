package mahasiswa_magang

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response intermoni.Response
	mahasiswa_magang intermoni.MahasiswaMagang
)

func GCFHandlerApplyMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	idmagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	err = ApplyMagang(idmagang, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Apply Magang"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerBatalApply(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	idmahasiswamagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	err = BatalApply(idmahasiswamagang, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Batal Apply"
	return intermoni.GCFReturnStruct(Response)

}

func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	idmahasiswamagang, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mahasiswa_magang)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role == "mitra" {
		if mahasiswa_magang.SeleksiBerkas != 0 {
			err = SeleksiBerkasMahasiswaMagangByMitra(idmahasiswamagang, user_login.Id, conn, mahasiswa_magang)
			if err != nil {
				Response.Message = err.Error()
				return intermoni.GCFReturnStruct(Response)
			}
			Response.Status = true
			Response.Message = "Berhasil Seleksi"
			return intermoni.GCFReturnStruct(Response)
		}
		err = SeleksiWawancaraMahasiswaMagangByMitra(idmahasiswamagang, user_login.Id, conn, mahasiswa_magang)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		Response.Status = true
		Response.Message = "Berhasil Seleksi"
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role == "mahasiswa" {
		err = UpdateStatusMahasiswaMagang(idmahasiswamagang, user_login.Id, conn, mahasiswa_magang)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		Response.Status = true
		Response.Message = "Berhasil Seleksi"
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerGetMahasiswaMagang(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id == "" {
		return GCFHandlerGetAllMahasiswaMagang(user_login, conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	mahasiswa_magang, err := intermoni.GetDetailMahasiswaMagangFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(mahasiswa_magang)
}

func GCFHandlerGetAllMahasiswaMagang(user_login intermoni.Payload, conn *mongo.Database) string {
	Response.Status = false
	//
	if user_login.Role == "admin" {
		mahasiswa_magang, err := GetMahasiswaMagangByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mahasiswa_magang)
	}
	if user_login.Role == "mitra" {
		mahasiswa_magang, err := GetMahasiswaMagangByMitra(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mahasiswa_magang)
	}
	if user_login.Role == "mahasiswa" {
		mahasiswa_magang, err := GetMahasiswaMagangByMahasiswa(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mahasiswa_magang)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}