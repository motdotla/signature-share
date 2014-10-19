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

type Params struct {
	DocumentUrl string `form:"document_url" json:"document_url"`
	SigningUrl  string `form:"signing_url" json:"signing_url"`
	SigningId   string `form:"signing_id" json:"signing_id"`
}

type SigningsCreateJson struct {
	Signings []Signing `json:"signings"`
}

type Signing struct {
	Id string `json:"id"`
}

type DocumentsShowJson struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	loadEnvs()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("assets"))

	m.Any("/", binding.Bind(Params{}), Index)

	m.Run()
}

func requestDocumentsShow(document_url string) DocumentsShowJson {
	res, err := http.Get(document_url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	var documents_show_json DocumentsShowJson
	json.NewDecoder(res.Body).Decode(&documents_show_json)

	return documents_show_json
}

func requestSigningsCreate(document_url string) SigningsCreateJson {
	signings_create_url := SIGNATURE_API_ROOT + "/api/v0/signings/create.json?document_url=" + document_url

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

	documents_show_json := requestDocumentsShow(document_url)
	status := documents_show_json.Documents[0].Status

	// if status is still processing than show a basic loading page
	if status == "processing" {
		mustaches := map[string]interface{}{"document_url": document_url}

		r.HTML(200, "processing", mustaches)
	} else {
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
}

func loadEnvs() {
	godotenv.Load()

	SIGNATURE_API_ROOT = os.Getenv("SIGNATURE_API_ROOT")
}
