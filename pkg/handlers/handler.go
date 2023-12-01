package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/S-F-Z/drosophila-stock-tracker/pkg/connectors"
	"github.com/S-F-Z/drosophila-stock-tracker/pkg/schema"
	"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func CreateUpdateHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var inputDmelTray *schema.DrosophilaTray
	var dbDmelTrays *schema.DrosophilaTray

	var vars = mux.Vars(r)

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "CreateHandler body data error %v"
		fmt.Fprintf(w, "%s", msg)
		return
	}

	conn.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &inputDmelTray)
	if errs != nil {
		msg := "EchoHandler could not unmarshal input data from json to schema %v"
		conn.Error(msg, errs)
		fmt.Fprintf(w, "%s", msg)
		return
	}

	// Open our jsonFile
	jsonFile, err := os.Open("tray_" + vars["id"] + ".json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// Initialize our Fly array

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'trays' which we defined above
	json.Unmarshal(byteValue, &dbDmelTrays)

	// tray = original, cr = input
	found := false
	var msg string
	for _, inptvial := range inputDmelTray.Vials {
		for _, dbvial := range dbDmelTrays.Vials {
			if dbvial.ID == inptvial.ID {
				dbvial = inptvial
				found = true
				msg = "Vial has been successfully updated."
			}
		}
		if !found {
			dbDmelTrays.Vials = append(dbDmelTrays.Vials, inptvial)
			msg = "Vial has been successfully added."
		}
	}

	// Prepare the tray data to be marshalled and written
	dataBytes, err := json.MarshalIndent(dbDmelTrays, "", " ")
	if err != nil {
		fmt.Println("Error Marshelling Tray:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	err = os.WriteFile("tray_1.json", dataBytes, fs.FileMode(jsonFile.Fd()))
	if err != nil {
		fmt.Println("Error when writing file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", msg)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var req *schema.Response

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", "error")
		return
	}

	conn.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &req)
	if errs != nil {
		msg := "CreateHandler could not unmarshal input data from json to schema %v"
		conn.Error(msg, errs)
		fmt.Fprintf(w, "%s", msg)
		return
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var trays []schema.DrosophilaTray
	var tray *schema.DrosophilaTray
	var list []byte

	files, err := FilePathWalkDir("json_files")
	if err != nil {
	}

	for _, file := range files {
		jsonFile, err := os.Open(file)
		if err != nil {

		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		if err := json.Unmarshal([]byte(byteValue), &tray); err != nil {
			// always check errors
		}

		trays = append(trays, *tray)

		list, err = json.MarshalIndent(trays, "", " ")
		if err != nil {
		}

	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", list)
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "%s", "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsFormat - utility function
func responseFormat(code int, status string, w http.ResponseWriter, payload schema.DrosophilaTray) string {
	response := `{"Code":"` + strconv.Itoa(code) + `", "Status": "` + status + `", "payload":"` + fmt.Sprintf(payload.Vials[0].ID) + `"}`
	w.WriteHeader(code)
	return response
}
