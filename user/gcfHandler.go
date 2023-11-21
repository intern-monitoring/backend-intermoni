package user

import (
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Response intermoni.Response
)

func Get(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role != "admin" {
		Response.Message = "Maneh bukan admin"
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id == "" {
		return GCFHandlerGetAllUserByAdmin(conn)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	user, err := intermoni.GetUserFromID(idparam, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user.Role == "mahasiswa" {
		mahasiswa, err := intermoni.GetMahasiswaFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mahasiswa)
	}
	if user.Role == "mitra" {
		mitra, err := intermoni.GetMitraFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mitra)
	}
	if user.Role == "pembimbing" {
		mitra, err := intermoni.GetPembimbingFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mitra)
	}
	if user.Role == "mentor" {
		mitra, err := intermoni.GetMentorFromAkun(user.ID, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mitra)
	}
	if user.Role == "admin" {
		mitra, err := intermoni.GetUserFromID(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mitra)
	}
	//
	Response.Message = "Tidak ada data"
	return intermoni.GCFReturnStruct(Response)
}

func GCFHandlerGetAllUserByAdmin(conn *mongo.Database) string {
	Response.Status = false
	//
	data, err := GetAllUser(conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	return intermoni.GCFReturnStruct(data)
}