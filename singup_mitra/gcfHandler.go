package singup_mitra

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
)

var (
	Response intermoni.Response
	mitra intermoni.Mitra
)

func Post(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&mitra)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = SignUpMitra(conn, mitra)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Halo " + mitra.Nama
	return intermoni.GCFReturnStruct(Response)
}