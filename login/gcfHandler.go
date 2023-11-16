package login

import (
	"encoding/json"
	"net/http"
	"os"

	intermoni "github.com/intern-monitoring/backend-intermoni"
)

var (
	Credential intermoni.Credential
	user intermoni.User
)

func Post(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := intermoni.MongoConnect(MONGOCONNSTRINGENV, dbname)
	Credential.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Credential.Message = "error parsing application/json: " + err.Error()
		return intermoni.GCFReturnStruct(Credential)
	}
	user, err := LogIn(conn, user)
	if err != nil {
		Credential.Message = err.Error()
		return intermoni.GCFReturnStruct(Credential)
	}
	tokenstring, err := intermoni.Encode(user.ID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Credential.Message = "Gagal Encode Token : " + err.Error()
		return intermoni.GCFReturnStruct(Credential)
	}
	//
	Credential.Message = "Selamat Datang " + user.Email
	Credential.Token = tokenstring
	Credential.Role = user.Role
	Credential.Status = true
	return intermoni.GCFReturnStruct(Credential)
}