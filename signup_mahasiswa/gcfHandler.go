package signup_mahasiswa

import (
	"encoding/json"
	"net/http"

	intermoni "github.com/intern-monitoring/backend-intermoni"
)

var (
	Response intermoni.Response
	mahasiswa intermoni.Mahasiswa
)

func Post(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	err = SignUpMahasiswa(conn, mahasiswa)
	if err != nil {
		Response.Message = err.Error()
		return intermoni.GCFReturnStruct(Response)
	}
	//
	Response.Status = true
	Response.Message = "Halo " + mahasiswa.NamaLengkap
	return intermoni.GCFReturnStruct(Response)
}