package seleksi

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Response intermoni.Response
	mahasiswa_magang intermoni.MahasiswaMagang
)

func Put(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	user_login, err := intermoni.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		Response.Message = err.Error()
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
	err = SeleksiBerkasMahasiswaMagangByMitra(idmahasiswamagang, user_login.Id, conn, mahasiswa_magang)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil Seleksi"
	return intermoni.GCFReturnStruct(Response)
}