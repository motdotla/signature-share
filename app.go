package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"net/http"
)

type Params struct {
	Url string `form:"url" json:"url"`
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Any("/", binding.Bind(Params{}), Index)

	m.Run()
}

func Index(params Params, req *http.Request, r render.Render) {
	_url := params.Url
	mustaches := map[string]interface{}{"url": _url}
	r.HTML(200, "index", mustaches)
}
