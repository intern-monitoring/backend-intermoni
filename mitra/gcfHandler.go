package mitra

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response intermoni.Response
	mitra    intermoni.Mitra
)

func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role == "admin" {
		return GCFHandlerUpdateByAdmin(user_login.Id, conn, r)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
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
	err = json.NewDecoder(r.Body).Decode(&mitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = UpdateMitra(idparam, user_login.Id, conn, mitra)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mitra"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerUpdateByAdmin(iduser primitive.ObjectID, conn *mongo.Database, r *http.Request) string {
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
	err = json.NewDecoder(r.Body).Decode(&mitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = ConfirmMouMitraByAdmin(idparam, conn, mitra)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mitra"
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
	if user_login.Role == "mitra" {
		return GetMitraByMitra(user_login.Id, conn)
	}
	if user_login.Role == "admin" {
		return GetMitraByAdmin(conn, r)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}

func GetMitraByMitra(iduser primitive.ObjectID, conn *mongo.Database) string {
	Response.Status = false
	//
	mitra, err := intermoni.GetMitraFromAkun(iduser, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(mitra)
}

func GetMitraByAdmin(conn *mongo.Database, r *http.Request) string {
	Response.Status = false
	//
	id := intermoni.GetID(r)
	if id == "" {
		mahasiswa, err := GetAllMitraByAdmin(conn)
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
	mahasiswa, err := GetMitraFromIDByAdmin(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(mahasiswa)
}