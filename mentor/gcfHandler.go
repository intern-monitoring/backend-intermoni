package mentor

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Response intermoni.Response
	mentor intermoni.Mentor
)

func Post(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	if user_login.Role != "mitra" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mentor)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = AddMentorByMitra(user_login.Id, conn, mentor)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Apply Magang"
	return intermoni.GCFReturnStruct(Response)
}

func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
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
	idmentor, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&mentor)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = UpdateMentor(idmentor, user_login.Id, conn, mentor)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Mentor"
	return intermoni.GCFReturnStruct(Response)
}

func Get(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if user_login.Role == "mitra" {
		if id == "" {
			mentor, err := GetAllMentorByMitra(user_login.Id, conn)
			if err != nil {
				Response.Message = err.Error()
				return intermoni.GCFReturnStruct(Response)
			}
			return intermoni.GCFReturnStruct(mentor)
		}
		idparam, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			Response.Message = "Invalid id parameter"
			return intermoni.GCFReturnStruct(Response)
		}
		mentor, err = intermoni.GetMentorFromID(idparam, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mentor)
	}
	if user_login.Role == "admin" {
		if id != "" {
			idparam, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				Response.Message = "Invalid id parameter"
				return intermoni.GCFReturnStruct(Response)
			}
			mentor, err = GetMentorFromIDByAdmin(idparam, conn)
			if err != nil {
				Response.Message = err.Error()
				return intermoni.GCFReturnStruct(Response)
			}
			return intermoni.GCFReturnStruct(mentor)
		}
		mentor, err := GetAllMentorByAdmin(conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mentor)
	}
	if user_login.Role == "mentor" {
		mentor, err := intermoni.GetMentorFromAkun(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(mentor)
	}
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}