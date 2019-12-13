// request.go kee > 2019/12/12

package request

import (
	"github.com/kataras/iris"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	ctx    iris.Context
	Method string
	Header http.Header
	Params RequestData
}

type RequestData struct {
	data map[string]interface{}
}

func New(ctx iris.Context) *Request {
	params := getParams(ctx)
	return &Request{
		ctx:    ctx,
		Header: ctx.Request().Header,
		Method: ctx.Request().Method,
		Params: params,
	}
}

func getParams(ctx iris.Context) RequestData {
	var params = make(map[string]interface{})

	switch ctx.Request().Header.Get("Content-Type") {
	case "application/json":
		ctx.ReadJSON(&params)
		break
	case "application/x-yaml", "text/yaml", "text/x-yaml":
		if body, err := GetBody(ctx.Request()); err == nil {
			yaml.Unmarshal(body, &params)
			pp := make(map[string]interface{})
			yaml.Unmarshal([]byte(`email:112233`), &pp)
			log.Println(&params, string(body), pp)
		}
		break
	case "application/xml", "text/xml":
		ctx.ReadXML(&params)
		break
	}

	forms := ctx.FormValues()
	for k, v := range forms {
		if _, ok := params[k]; !ok {
			params[k] = v
		}
	}
	return RequestData{params}
}

func (req *Request) Get(name string) interface{} {
	return req.Params.Get(name)
}

func (req *Request) GetBody() ([]byte, error) {
	return GetBody(req.ctx.Request())
}

func (p *RequestData) Get(name string) (val interface{}) {
	return p.data[name]
}

func GetBody(r *http.Request) ([]byte, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, err
}
