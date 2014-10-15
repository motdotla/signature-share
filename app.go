package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"net/http"
)

type Params struct {
	DocumentUrl string `form:"document_url" json:"document_url"`
	SigningUrl  string `form:"signing_url" json:"signing_url"`
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Any("/", binding.Bind(Params{}), Index)

	m.Run()
}

func Index(params Params, req *http.Request, r render.Render) {
	document_url := params.DocumentUrl
	signing_url := params.DocumentUrl
	mustaches := map[string]interface{}{"document_url": document_url, "signing_url": signing_url}
	r.HTML(200, "index", mustaches)
}
