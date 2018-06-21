package mfe

import (
	"encoding/json"
	"fmt"
	"sort"

	fasthttp "github.com/valyala/fasthttp"
)

// InputData -- Входные параметры запроса
type InputData struct {
	path    string
	method  string
	get     map[string]string
	post    map[string]string
	json    Variant
	cookie  map[string]string
	headers map[string]string
}

var rDir map[string]func(ctx *fasthttp.RequestCtx)
var rMap map[int]map[string]func(ctx *fasthttp.RequestCtx)
var lMap map[int]struct{}
var lS []int

var serverHeader string

var defRoute func(ctx *fasthttp.RequestCtx)

// ServerHeaderSet -- Установить значение в хедер сервер
func ServerHeaderSet(server string) {
	serverHeader = server
}

// AddDirectRoute -- Добавить роут точное имя
func AddDirectRoute(uri string, f func(ctx *fasthttp.RequestCtx)) {
	uri = "/" + uri
	if rDir == nil {
		rDir = make(map[string]func(ctx *fasthttp.RequestCtx))
	}
	rDir[uri] = f
}

// AddRoute -- Добавить роут
func AddRoute(uri string, f func(ctx *fasthttp.RequestCtx)) {
	uri = "/" + uri
	if rMap == nil {
		rMap = make(map[int]map[string]func(ctx *fasthttp.RequestCtx))
		lMap = make(map[int]struct{})
	}

	l := len(uri)
	if _, ok := lMap[l]; !ok {
		lMap[l] = struct{}{}
		rMap[l] = map[string]func(ctx *fasthttp.RequestCtx){}
		lS = append(lS, l)
		sort.Slice(lS, func(i, j int) bool { return lS[i] > lS[j] })
	}

	rMap[l][uri] = f
}

// AddDefaultRoute -- Добавить роут по умолчанию
func AddDefaultRoute(f func(ctx *fasthttp.RequestCtx)) {
	defRoute = f
}

func searchRoute(uri string) (fn func(ctx *fasthttp.RequestCtx)) {

	if f, ok := rDir[uri]; ok {
		return f
	}

	l := len(uri)

	for _, v := range lS {
		if v <= l {
			s := uri[:v]
			m := rMap[v]
			if f, ok := m[s]; ok {
				return f
			}
		}
	}

	return defRoute
}

// ListenAndServe Запуск fastHttp сервиса
func ListenAndServe(listenAddr string, defaultHandler func(ctx *fasthttp.RequestCtx)) {
	AddDefaultRoute(defaultHandler)
	fasthttp.ListenAndServe(listenAddr, func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.SetServer(StringGetNotEmpty(serverHeader, "mfe"))

		f := searchRoute(string(ctx.Path()))

		f(ctx)
	})
}

// DisplayInputHandler отображает всё что пришло на вход
func DisplayInputHandler(ctx *fasthttp.RequestCtx) {

	fmt.Fprintln(ctx, string("path: "))
	fmt.Fprintln(ctx, string(ctx.Path()))
	fmt.Fprintln(ctx, "<br>method: ")
	fmt.Fprintln(ctx, string(ctx.Method()))
	fmt.Fprintln(ctx, "<br>post args: ")
	fmt.Fprintln(ctx, ctx.PostArgs())
	fmt.Fprintln(ctx, "<br>post body: ")
	fmt.Fprintln(ctx, string(ctx.PostBody()))
	fmt.Fprintln(ctx, "<br>query args: ")

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		fmt.Fprintln(ctx, "<br>")
		fmt.Fprint(ctx, string(key))
		fmt.Fprint(ctx, " : ")
		fmt.Fprintln(ctx, string(value))
	})

	fmt.Fprintln(ctx, "<br>headrs:")
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		fmt.Fprintln(ctx, "<br>")
		fmt.Fprint(ctx, string(key))
		fmt.Fprint(ctx, " : ")
		fmt.Fprintln(ctx, string(value))
	})

	fmt.Fprintln(ctx, "<br>cookie:")
	ctx.Request.Header.VisitAllCookie(func(key, value []byte) {
		fmt.Fprintln(ctx, "<br>")
		fmt.Fprint(ctx, string(key))
		fmt.Fprint(ctx, " : ")
		fmt.Fprintln(ctx, string(value))
	})

	ctx.Response.Header.SetContentType("text/html;charset=utf-8")
}

// InputDataGet получить все значения из
func InputDataGet(ctx *fasthttp.RequestCtx) (id InputData) {
	id = InputData{}
	id.get = map[string]string{}
	id.post = map[string]string{}
	id.headers = map[string]string{}
	id.cookie = map[string]string{}

	id.path = string(ctx.Path())
	id.method = string(ctx.Method())

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		id.get[string(key)] = string(value)
	})

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		id.headers[string(key)] = string(value)
	})

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		id.post[string(key)] = string(value)
	})

	ctx.Request.Header.VisitAllCookie(func(key, value []byte) {
		id.cookie[string(key)] = string(value)
	})

	var vv Variant
	if err := json.Unmarshal(ctx.PostBody(), &vv); err == nil {
		id.json = vv
	} else {
		id.json = VariantNewNull()
	}

	return id
}

// EmptyHandler пустой ответ
func EmptyHandler(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.SetContentType("text/html;charset=utf-8")
}

// VariantHandler выдача variant
func VariantHandler(ctx *fasthttp.RequestCtx, v Variant, err error, errMsg Variant) {
	ctx.Response.Header.SetContentType("application/json;charset=utf-8")

	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	if !errMsg.IsNull() {
		ctx.Response.SetStatusCode(400)
		fmt.Fprint(ctx, v)
		return
	}

	fmt.Fprint(ctx, v)
}

// Path Get
func (id InputData) Path() string {
	return id.path
}

// Method Get
func (id InputData) Method() string {
	return id.method
}

// Get Get
func (id InputData) Get() map[string]string {
	return id.get
}

// Post Get
func (id InputData) Post() map[string]string {
	return id.post
}

// Cookie Get
func (id InputData) Cookie() map[string]string {
	return id.cookie
}

// Headers Get
func (id InputData) Headers() map[string]string {
	return id.headers
}

// JSONData Get
func (id InputData) JSONData() Variant {
	return id.json
}

// ParamGet Get parametr
func (id InputData) ParamGet(n string) (s string) {
	if v, isOk := id.json.GetElement(n); isOk {
		return v.String()
	}

	if v, isOk := id.post[n]; isOk {
		return v
	}

	if v, isOk := id.get[n]; isOk {
		return v
	}

	return ""
}

// HTTPCallFull do http request with fasthttp
func HTTPCallFull(url string, method string, body []byte, headers map[string]string, cookies map[string]string) (statusCode int, bodyBytes []byte, headersOut map[string]string, cookiesOut map[string]string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if body != nil && len(body) > 0 {
		req.SetBody(body)
	}

	if headers != nil {
		for k, v := range headers {
			if k == "User-Agent" {
				req.Header.SetUserAgent(v)
			} else {
				req.Header.Add(k, v)
			}
		}

	}

	if cookies != nil {
		for k, v := range cookies {
			req.Header.SetCookie(k, v)
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	if err != nil {
		return
	}

	bodyBytes = resp.Body()

	statusCode = resp.Header.StatusCode()

	headersOut = map[string]string{}

	resp.Header.VisitAll(func(key, value []byte) {
		headersOut[string(key)] = string(value)
	})

	cookiesOut = map[string]string{}

	resp.Header.VisitAllCookie(func(key, value []byte) {
		cookiesOut[string(key)] = string(value)
	})

	return
}