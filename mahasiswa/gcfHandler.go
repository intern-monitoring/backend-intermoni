package mahasiswa

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response intermoni.Response
	mahasiswa intermoni.Mahasiswa
)

func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	id := intermoni.GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role == "mahasiswa" {
		return GCFHandlerUpdateByMahasiswa(idparam, user_login.Id, mahasiswa, conn, r)
	}
	if user_login.Role == "admin" {
		return GCFHandlerUpdateByAdmin(idparam, mahasiswa, conn, r)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerUpdateByMahasiswa(idparam, iduser primitive.ObjectID,  mahasiswa intermoni.Mahasiswa, conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	err := UpdateMahasiswa(idparam, iduser, conn, mahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mahasiswa"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerUpdateByAdmin(idparam primitive.ObjectID, mahasiswa intermoni.Mahasiswa, conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	err := SeleksiMahasiswaByAdmin(idparam, conn, mahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Seleksi Mahasiswa"
	return intermoni.GCFReturnStruct(Response)
}

func Get(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role == "mahasiswa" {
		return GCFHandlerGetMahasiswaByMahasiswa(user_login.Id, conn)
	}
	if user_login.Role == "admin" {
		return GCFHandlerGetMahasiswaByAdmin(conn, r)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerGetMahasiswaByMahasiswa(iduser primitive.ObjectID, conn *mongo.Database) string {
	Response.Status = false
	//
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(mahasiswa)
}

func GCFHandlerGetMahasiswaByAdmin(conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	id := intermoni.GetID(r)
	if id == "" {
		mahasiswa, err := GetAllMahasiswaByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mahasiswa)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	mahasiswa, err := GetMahasiswaFromIDByAdmin(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(mahasiswa)
}