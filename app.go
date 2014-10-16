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

type SigningsCreateJson struct {
	Signings []Signing `json:"signings"`
}

type Signing struct {
	Id string `json:"id"`
}

func main() {
	loadEnvs()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(CrossDomain())

	m.Any("/", binding.Bind(Params{}), Index)

	m.Run()
}

func requestSigningsCreate(document_url string) SigningsCreateJson {
	signings_create_url := SIGNATURE_API_ROOT + "/api/v0/signings/create.json?document_url=" + document_url

	log.Println(signings_create_url)

	res, err := http.Get(signings_create_url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	var signings_create_json SigningsCreateJson
	json.NewDecoder(res.Body).Decode(&signings_create_json)
	return signings_create_json
}

func Index(params Params, req *http.Request, r render.Render) {
	document_url := params.DocumentUrl
	signing_url := params.SigningUrl

	if signing_url != "" {
		mustaches := map[string]interface{}{"document_url": document_url, "signing_url": signing_url}
		r.HTML(200, "index", mustaches)
	} else {
		signings_create_json := requestSigningsCreate(document_url)
		signing_id := signings_create_json.Signings[0].Id
		created_signing_url := SIGNATURE_API_ROOT + "/api/v0/signings/" + signing_id + ".json"

		r.Redirect("/?document_url=" + document_url + "&signing_url=" + created_signing_url)
	}
}

func loadEnvs() {
	godotenv.Load()

	SIGNATURE_API_ROOT = os.Getenv("SIGNATURE_API_ROOT")
}
