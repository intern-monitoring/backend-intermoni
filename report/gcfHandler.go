package report

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Response intermoni.Response
	report intermoni.Report
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
	if user_login.Role != "mahasiswa" {
		Response.Message = "Maneh tidak memiliki akses"
		return intermoni.GCFReturnStruct(Response)
	}
	err = json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = TambahReportByMahasiswa(user_login.Id, conn, report)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Menambahkan Report"
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
	err = json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = UpdateReportByMahasiswa(idparam, user_login.Id, conn, report)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Update Report"
	return intermoni.GCFReturnStruct(Response)
}

func Delete(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
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
	err = DeleteReportByMahasiswa(idparam, user_login.Id, conn)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Berhasil Delete Report"
	return intermoni.GCFReturnStruct(Response)
}

func Get(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = true
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = "Gagal Decode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	id := intermoni.GetID(r)
	if id != "" {
		idparam, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			Response.Message = "Invalid id parameter"
			return intermoni.GCFReturnStruct(Response)
		}
		report, err := GetReportByID(idparam, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(report)
	}
	if user_login.Role == "mahasiswa" {
		report, err := GetAllReportByMahasiswa(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(report)
	}
	if user_login.Role == "mentor" || user_login.Role == "pembimbing" {
		report, err := GetAllReportByPenerima(user_login.Id, conn)
		if err != nil {
			Response.Message = err.Error()
			return intermoni.GCFReturnStruct(Response)
		}
		return intermoni.GCFReturnStruct(report)
	}
	//
	Response.Message = "Maneh tidak memiliki akses"
	return intermoni.GCFReturnStruct(Response)
}