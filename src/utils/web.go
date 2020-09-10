package utils

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/i-PUSH/RPi-Cast/bindata"
	"github.com/kabukky/httpscerts"
)

// GetURLParam and decode params... need to do
func GetURLParam(r *http.Request, k string) string {
	return r.URL.Query()[k][0]
}

// ExecTemplate with name and data
func ExecTemplate(w http.ResponseWriter, fileName string, data interface{}) {
	asset, err1 := bindata.Asset(fileName)
	tmpl, err2 := template.New(fileName).Parse(string(asset))
	if err1 != nil || err2 != nil {
		log.Println(err1, " | ", err2)
	}
	LogErr(tmpl.Execute(w, data))
}

// ReadJSON return decoded struct
func ReadJSON(req *http.Request, v interface{}) {
	decoder := json.NewDecoder(req.Body)
	LogErr(decoder.Decode(&v))
	defer req.Body.Close()
}

// WriteJSON create json object for response
func WriteJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// CheckHTTPS and generate https certs
func CheckHTTPS(addr string) {
	// Check if the cert files are available.
	err := httpscerts.Check("cert.pem", "key.pem")
	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", addr)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
}
