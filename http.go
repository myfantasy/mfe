package mfe

import (
	"fmt"
	"sort"

	fasthttp "github.com/valyala/fasthttp"
)

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
	rDir[uri] = f
}

// AddRoute -- Добавить роут
func AddRoute(uri string, f func(ctx *fasthttp.RequestCtx)) {
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

	ctx.Response.Header.SetContentType("text/html;charset=utf-8")
}
