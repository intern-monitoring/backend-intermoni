package pembimbing

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Response intermoni.Response
	pembimbing intermoni.Pembimbing
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
	if user_login.Role != "admin" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&pembimbing)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = AddPembimbingByAdmin(conn, pembimbing)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Add Pembimbing"
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
	if user_login.Role != "pembimbing" {
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
	err = json.NewDecoder(r.Body).Decode(&pembimbing)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = UpdatePembimbing(idparam, user_login.Id, conn, pembimbing)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Pembimbing"
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
	if user_login.Role == "admin" {
		if id == "" {
			pembimbing, err := GetAllPembimbingByAdmin(conn)
			if err != nil {
				Response.Message = err.Error()
				return intermoni.GCFReturnStruct(Response)
			}
			return intermoni.GCFReturnStruct(pembimbing)
		}
		idparam, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		pembimbing, err := intermoni.GetPembimbingFromID(idparam, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(pembimbing)
	}
	if user_login.Role == "pembimbing" {
		pembimbing, err := intermoni.GetPembimbingFromAkun(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(pembimbing)
	}
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}