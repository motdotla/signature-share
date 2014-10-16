package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/binding"
	"log"
	"net/http"
	"os"
)

var (
	SIGNATURE_API_ROOT string
)

func CrossDomain() martini.Handler {
	return func(res http.ResponseWriter) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
}

type Params struct {
	DocumentUrl string `form:"document_url" json:"document_url"`
	SigningUrl  string `form:"signing_url" json:"signing_url"`
}

func main() {
	loadEnvs()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(CrossDomain())

	m.Any("/", binding.Bind(Params{}), Index)

	m.Run()
}

func requestSigningsCreate(document_url string) interface{} {
	signings_create_url := SIGNATURE_API_ROOT + "/api/v0/signings/create.json?document_url=" + document_url

	log.Println(signings_create_url)

	res, err := http.Get(signings_create_url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	return data
}

func Index(params Params, req *http.Request, r render.Render) {
	document_url := params.DocumentUrl
	signing_url := params.SigningUrl

	if signing_url != "" {
		mustaches := map[string]interface{}{"document_url": document_url, "signing_url": signing_url}
		r.HTML(200, "index", mustaches)
	} else {
		json_response := requestSigningsCreate(document_url)

		// yuck
		signings := json_response.(map[string]interface{})["signings"].([]interface{})
		signing := signings[0].(map[string]interface{})
		signing_id := signing["id"].(string)
		created_signing_url := SIGNATURE_API_ROOT + "/api/v0/signings/" + signing_id + ".json"

		r.Redirect("/?document_url=" + document_url + "&signing_url=" + created_signing_url)
	}
}

func loadEnvs() {
	godotenv.Load()

	SIGNATURE_API_ROOT = os.Getenv("SIGNATURE_API_ROOT")
}
